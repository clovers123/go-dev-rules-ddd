package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ddd-bootstrap/internal/config"
)

const initDBInfoPath = ".agent/.init-db-info.json"

type initDBInfo struct {
	Mode                   string `json:"mode"`
	DB                     string `json:"db"`
	Host                   string `json:"host"`
	Port                   int    `json:"port"`
	User                   string `json:"user"`
	Password               string `json:"-"`
	PasswordStoredInConfig bool   `json:"passwordStoredInConfig"`
	Database               string `json:"database"`
	SSLMode                string `json:"sslMode"`
	ContainerName          string `json:"containerName,omitempty"`
	ContainerImage         string `json:"containerImage,omitempty"`
	CreatedAt              string `json:"createdAt"`
}

func prepareInitDBInfo(opts *config.InitOptions) (*initDBInfo, error) {
	if !opts.AutoGormGen && !opts.UseExistingDB {
		return nil, nil
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)
	if opts.DryRun {
		createdAt = "dry-run"
	}
	info := &initDBInfo{
		DB:                     opts.DB,
		Host:                   opts.DBHost,
		Port:                   opts.DBPort,
		User:                   opts.DBUser,
		Password:               opts.DBPassword,
		PasswordStoredInConfig: true,
		Database:               opts.DBName,
		SSLMode:                opts.DBSSLMode,
		CreatedAt:              createdAt,
	}

	if opts.UseExistingDB {
		info.Mode = "use-existing-db"
		return info, nil
	}

	port := opts.DBPort
	if !opts.DryRun {
		var err error
		port, err = randomLocalPort()
		if err != nil {
			return nil, err
		}
	}
	opts.DBHost = "127.0.0.1"
	opts.DBPort = port

	info.Mode = "auto-gormgen"
	info.Host = opts.DBHost
	info.Port = opts.DBPort
	info.ContainerImage = "postgres:16-alpine"
	if opts.DryRun {
		info.ContainerName = "ddd-bootstrap-" + normalizeContainerPart(opts.ProjectName) + "-dry-run"
	} else {
		info.ContainerName = "ddd-bootstrap-" + normalizeContainerPart(opts.ProjectName) + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}
	return info, nil
}

func renderInitDBInfo(info *initDBInfo) (plannedFile, error) {
	content, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return plannedFile{}, fmt.Errorf("render init db info: %w", err)
	}
	return plannedFile{target: initDBInfoPath, content: string(content) + "\n"}, nil
}

func activateInitDB(out io.Writer, root string, info *initDBInfo) error {
	if info.Mode != "auto-gormgen" {
		return nil
	}
	if err := ensureDockerAvailable(); err != nil {
		return err
	}

	fmt.Fprintf(out, "RUN docker start postgres container %s on %s:%d\n", info.ContainerName, info.Host, info.Port)
	cmd := exec.Command(
		"docker", "run", "-d",
		"--name", info.ContainerName,
		"-e", "POSTGRES_USER="+info.User,
		"-e", "POSTGRES_PASSWORD="+info.Password,
		"-e", "POSTGRES_DB="+info.Database,
		"-p", fmt.Sprintf("%s:%d:5432", info.Host, info.Port),
		info.ContainerImage,
	)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("start postgres container: %w\n%s", err, strings.TrimSpace(string(output)))
	}
	return waitForPostgresContainer(info)
}

func runGormGenIfConfigured(out io.Writer, root string) error {
	info, err := readInitDBInfo(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.Mode != "auto-gormgen" && info.Mode != "use-existing-db" {
		return nil
	}
	if info.Mode == "auto-gormgen" {
		if err := ensurePostgresContainer(out, info); err != nil {
			return err
		}
	}
	return runGormGen(out, root)
}

var runGormGen = func(out io.Writer, root string) error {
	fmt.Fprintln(out, "RUN go run -mod=mod ./cmd/gorm-gen")
	cmd := exec.Command("go", "run", "-mod=mod", "./cmd/gorm-gen")
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Fprint(out, string(output))
	}
	if err != nil {
		return fmt.Errorf("run gorm gen: %w", err)
	}
	return nil
}

var runGoModDownload = func(out io.Writer, root string) error {
	fmt.Fprintln(out, "RUN go mod download all")
	cmd := exec.Command("go", "mod", "download", "all")
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Fprint(out, string(output))
	}
	if err != nil {
		return fmt.Errorf("run go mod download: %w", err)
	}
	return nil
}

var runGoModTidy = func(out io.Writer, root string) error {
	fmt.Fprintln(out, "RUN go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Fprint(out, string(output))
	}
	if err != nil {
		return fmt.Errorf("run go mod tidy: %w", err)
	}
	return nil
}

func readInitDBInfo(root string) (*initDBInfo, error) {
	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(initDBInfoPath)))
	if err != nil {
		return nil, err
	}
	var info initDBInfo
	if err := json.Unmarshal(content, &info); err != nil {
		return nil, fmt.Errorf("read %s: %w", initDBInfoPath, err)
	}
	return &info, nil
}

func ensurePostgresContainer(out io.Writer, info *initDBInfo) error {
	if info.ContainerName == "" {
		return fmt.Errorf("auto-gormgen database info is missing containerName")
	}
	if err := ensureDockerAvailable(); err != nil {
		return err
	}
	inspect := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", info.ContainerName)
	output, err := inspect.CombinedOutput()
	if err == nil && strings.TrimSpace(string(output)) == "true" {
		return nil
	}
	fmt.Fprintf(out, "RUN docker start %s\n", info.ContainerName)
	start := exec.Command("docker", "start", info.ContainerName)
	if startOutput, startErr := start.CombinedOutput(); startErr != nil {
		return fmt.Errorf("start existing postgres container: %w\n%s", startErr, strings.TrimSpace(string(startOutput)))
	}
	return waitForPostgresContainer(info)
}

func ensureDockerAvailable() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is unavailable; skip --auto-gormgen because Docker is not installed or not on PATH")
	}
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker is unavailable; skip --auto-gormgen because the Docker daemon cannot be reached: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func waitForPostgresContainer(info *initDBInfo) error {
	var lastOutput []byte
	for i := 0; i < 30; i++ {
		cmd := exec.Command("docker", "exec", info.ContainerName, "pg_isready", "-U", info.User, "-d", info.Database)
		output, err := cmd.CombinedOutput()
		if err == nil {
			return nil
		}
		lastOutput = output
		time.Sleep(time.Second)
	}
	return fmt.Errorf("postgres container %s did not become ready: %s", info.ContainerName, strings.TrimSpace(string(lastOutput)))
}

func randomLocalPort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("allocate local port: %w", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func normalizeContainerPart(value string) string {
	value = strings.ToLower(value)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	value = strings.Trim(re.ReplaceAllString(value, "-"), "-")
	if value == "" {
		return "project"
	}
	if len(value) > 32 {
		return value[:32]
	}
	return value
}
