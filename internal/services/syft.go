package services

import (
	"fmt"
	"os"
	"os/exec"
)

func ExtractDependenciesManual(repo string) ([]Package, error) {
	dir, err := getRepo(repo)
	if err != nil {
		return nil, err
	}

	sbomPath := dir + "/sbom.json"

	err = generateSbom(dir, sbomPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(sbomPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ExtractDependencies(f)
}

func getRepo(repo string) (string, error) {
	dir := "/tmp/repos/" + repo
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		fmt.Println("Cloning repo...")
		cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/"+repo+".git", dir)
		err := cmd.Run()
		if err != nil {
			return "", err
		}
	} else {
		fmt.Println("Repo exists... running git fetch")
		fetchCmd := exec.Command("git", "-C", dir, "fetch", "--depth", "1")
		if output, err := fetchCmd.CombinedOutput(); err != nil {
			return fmt.Sprintf("fetch failed: %v\n%s", err, output), err
		}

		resetCmd := exec.Command("git", "-C", dir, "reset", "--hard", "origin/HEAD")
		if output, err := resetCmd.CombinedOutput(); err != nil {
			return fmt.Sprintf("reset failed: %v\n%s", err, output), nil
		}

	}
	return dir, nil
}

func generateSbom(repoDir, sbomPath string) error {
	cmd := exec.Command(
		"syft",
		repoDir,
		"-o",
		"spdx-json="+sbomPath,
	)

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("SBOM generated:", sbomPath)
	return nil
}
