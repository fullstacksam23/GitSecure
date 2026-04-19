package sbom

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fullstacksam23/GitSecure/internal/core"
)

func ExtractDependenciesManual(repo string) ([]core.Package, []byte, error) {
	dir, err := getRepo(repo)
	if err != nil {
		return nil, nil, err
	}

	sbomPath := dir + "/sbom.json"

	err = generateSbom(dir, sbomPath)
	if err != nil {
		return nil, nil, err
	}

	data, err := os.ReadFile(sbomPath)
	if err != nil {
		return nil, nil, err
	}

	pkgs, err := ExtractDependencies(data)
	if err != nil {
		return nil, nil, err
	}

	return pkgs, data, nil
}

func getRepo(repo string) (string, error) {
	dir := filepath.Join(os.TempDir(), "gitsecure", "repos", filepath.FromSlash(repo))
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return "", fmt.Errorf("failed to create repo cache dir: %w", err)
	}

	_, err := os.Stat(dir)

	//repo doesnt exist locally
	if os.IsNotExist(err) {
		startedAt := time.Now()
		log.Printf("cloning repo %s into %s", repo, dir)
		cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/"+repo+".git", dir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("git clone failed for %s after %s: %s", repo, time.Since(startedAt).Round(time.Second), string(output))
			return "", fmt.Errorf("git clone failed: %w", err)
		}
		log.Printf("clone finished for %s in %s", repo, time.Since(startedAt).Round(time.Second))
	} else {
		log.Printf("repo cache exists for %s, refreshing checkout", repo)
		fetchCmd := exec.Command("git", "-C", dir, "fetch", "--depth", "1")
		if output, err := fetchCmd.CombinedOutput(); err != nil {
			log.Println("fetch failed: ", err, string(output))
			return "", fmt.Errorf("git fetch failed: %w", err)
		}

		resetCmd := exec.Command("git", "-C", dir, "reset", "--hard", "origin/HEAD")
		if output, err := resetCmd.CombinedOutput(); err != nil {
			log.Println("reset failed: ", err, string(output))
			return "", fmt.Errorf("git reset failed: %w", err)
		}

	}
	return dir, nil
}

func generateSbom(repoDir, sbomPath string) error {
	startedAt := time.Now()
	log.Printf("generating sbom for %s", repoDir)
	cmd := exec.Command(
		"syft",
		repoDir,
		"-o",
		"spdx-json="+sbomPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("syft failed for %s after %s: %s", repoDir, time.Since(startedAt).Round(time.Second), string(output))
		return fmt.Errorf("syft failed: %w", err)
	}

	log.Printf("sbom generated at %s in %s", sbomPath, time.Since(startedAt).Round(time.Second))
	return nil
}
