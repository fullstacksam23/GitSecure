package services

import (
	"fmt"
	"os"
	"os/exec"
)

func ExtractDependenciesManual(repo string) ([]Package, error) {
	dir, err := cloneRepo(repo)
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
func cloneRepo(repo string) (string, error) {
	dir := "/tmp/repos/" + repo

	//remove dir if exists already to replace with updated contents
	os.RemoveAll(dir)

	cmd := exec.Command("git", "clone", "https://github.com/"+repo+".git", dir)
	err := cmd.Run()
	if err != nil {
		return "", err
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
