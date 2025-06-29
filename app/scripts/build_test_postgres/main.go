package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// This script creates a test postgres container and then publishes it to GitHub.
// This image will have all of the most current migrations and be seeded with test data.
// The image can be used to run integration tests locally or in GitHub CI/CD.
// This image is created this way so that the spin up time of the container is fast since all of the migrations and seed data have already been loaded.
// All data will be reset every time the container is ran.

func main() {
	// Get the working directory
	cwd := getCwd()
	log.Print("Working Directory: " + cwd)

	// Check if Docker Dameon is running
	// Docker Dameon needs to be started outside of this script.
	checkDockerRunning()

	// Load env vars
	err := godotenv.Load(".env.tst")
	if err != nil {
		log.Fatalf("Error: B2WRUD - Loading .env.inttest Error: %v", err)
	}

	// Get GitHub user and token
	githubUsername := os.Getenv("GIT_HUB_USERNAME")
	githubToken := os.Getenv("GIT_HUB_TOKEN")
	if githubUsername == "" || githubToken == "" {
		log.Fatal("GIT_HUB_USERNAME and GIT_HUB_TOKEN environment variables must be set")
	}

	// Create temp directory paths and Docker image tag
	tempDataDir := filepath.Join(cwd, "temp/temp_pgdata")
	tempMigrationsDir := filepath.Join(cwd, "temp/temp_migrations")
	imageTag := "ghcr.io/" + githubUsername + "/test-postgres:latest"

	// Clean old temp dirs
	_ = os.RemoveAll(tempDataDir)
	_ = os.RemoveAll(tempMigrationsDir)

	// Copy only .up.sql files from migrations directory into a temp init directory
	copyUpSQLFiles(filepath.Join(cwd, "migrations"), tempMigrationsDir)

	// Copy testing seed sql file to the temp init directory
	seedFile := filepath.Join(cwd, "seeds/testing/999999_seed_test_db.sql")
	copyFile(seedFile, filepath.Join(tempMigrationsDir, "999999_seed_test_db.sql"))

	// Run a postgres container to initialize a test db with all of the migrations and seed data.
	// This will place the new db in the tempDataDir where we can copy it into the final test-postgres image.
	// The db will be initialized with all the sql files located in tempMigrationsDir.
	run("docker", "run", "-d", "--name", "pg-init",
		"-e", "POSTGRES_PASSWORD=pass",
		"-v", tempDataDir+":/var/lib/postgresql/data",
		"-v", tempMigrationsDir+":/docker-entrypoint-initdb.d/",
		"postgres:latest")

	// Wait for initialization to complete
	log.Println("⏳ Waiting for Postgres to finish initialization...")

	// Run a bash shell command in the initialize db container to check when postgres is available and taking connections.
	// It will check in a loop waiting one second between checks using the user postgres.
	// This means the db has been created and initialized and we can move on.
	run("docker", "exec", "pg-init", "bash", "-c", "until pg_isready -U postgres; do sleep 1; done")

	// Give 5 seconds time for scripts to run
	// This may need to be adjusted based on the length of the sql scripts
	log.Println("✅ Postgres is ready, waiting a few seconds for SQL scripts...")
	run("sleep", "5")

	// Stop and remove the initialize db container
	run("docker", "stop", "pg-init")
	run("docker", "rm", "pg-init")

	// Build final Docker image that contains pre-seeded DB with all of the most current migrations.
	// See tests/Dockerfile. This is where the db is copied into the image.
	run("docker", "build",
		"-f", filepath.Join(cwd, "int_testing/Dockerfile"),
		"--build-arg", "POSTGRES_PASSWORD=pass",
		"-t", imageTag,
		".")

	// GitHub login using username and token from env vars
	dockerLoginToGitHub(githubUsername, githubToken)

	// Push the image to GitHub
	run("docker", "push", imageTag)

	// Clean up temp folder
	_ = os.RemoveAll(tempDataDir)
	_ = os.RemoveAll(tempMigrationsDir)

	log.Println("✅ Done. Image pushed:", imageTag)
}

// A wrapper function to execute os command line commands.
func run(cmd string, args ...string) {
	log.Printf("Running: %s %v\n", cmd, args)
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		log.Fatalf("Command failed: %s %v\nError: %v", cmd, args, err)
	}
}

// A function to log in to GitHub
func dockerLoginToGitHub(githubUsername, githubToken string) {
	cmd := exec.Command("docker", "login", "ghcr.io", "-u", githubUsername, "--password-stdin")
	cmd.Stdin = strings.NewReader(githubToken)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("GitHub Docker login failed: %v", err)
	}
}

// A function used to copy only "up" migration sql to the tempMigrationsDir.
func copyUpSQLFiles(srcDir, destDir string) {
	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatalf("Failed to read migrations dir: %v", err)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Fatalf("Failed to create temp migrations dir: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			src := filepath.Join(srcDir, file.Name())
			dest := filepath.Join(destDir, file.Name())
			copyFile(src, dest)
		}
	}
}

// A helper function to copy files
func copyFile(src, dest string) {
	in, err := os.Open(src)
	if err != nil {
		log.Fatalf("Failed to open source file %s: %v", src, err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		log.Fatalf("Failed to create directory for %s: %v", dest, err)
	}

	out, err := os.Create(dest)
	if err != nil {
		log.Fatalf("Failed to create destination file %s: %v", dest, err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		log.Fatalf("Failed to copy %s to %s: %v", src, dest, err)
	}
}

// A helper function to get the working directory path
func getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

// A helper function to check if the Docker Dameon is running
func checkDockerRunning() {
	cmd := exec.Command("docker", "info")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		log.Fatal("❌ Docker is not running. Please start Docker Desktop or the Docker daemon.")
	}
}
