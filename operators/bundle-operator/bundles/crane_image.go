package bundles

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var bundleLog = ctrl.Log.WithName("bundles")

func ExtractImageTo(repoSrc string, pathDir string) error {
	pathTarFile := generateFileTarPath(pathDir)
	f, err := openFile(pathTarFile)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", pathTarFile, err)
	}
	defer f.Close()

	// pull image
	var img v1.Image
	var options *[]crane.Option = &[]crane.Option{}
	img, err = crane.Pull(repoSrc, *options...)
	if err != nil {
		return fmt.Errorf("pulling %s: %w", repoSrc, err)
	}

	// mutate extract to pathDir.tar
	err = crane.Export(img, f)
	if err != nil {
		return fmt.Errorf("failed to export %s: %w", repoSrc, err)
	}

	// untar pathDir.tar
	err = extractTar(pathDir, pathTarFile)
	if err != nil {
		return fmt.Errorf("failed to extract tar %s: %w", pathTarFile, err)
	}

	return os.Remove(pathTarFile)
}

func generateFileTarPath(pathDir string) string {
	filePath := strings.TrimSuffix(pathDir, "/")
	return filePath + ".tar"
}

func openFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			return nil, err
		}
	}
	f, err := os.Create(path)
	return f, err
}

func extractTar(pathDir string, pathTarFile string) error {
	tarFile, err := os.Open(pathTarFile)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	var fileReader io.ReadCloser = tarFile
	tarBallReader := tar.NewReader(fileReader)
	fmt.Println("extractTar start")
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				fmt.Println("extractTar EOF")
				break
			}
			return err
		}

		// get the individual filename and extract to the current directory
		filename := header.Name
		target := filepath.Join(pathDir, filename)

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			bundleLog.Info(fmt.Sprintf("Creating directory : %s", target))
			err = os.MkdirAll(target, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				return err
			}

		case tar.TypeReg:
			// handle normal file
			bundleLog.Info(fmt.Sprintf("Untarring : %s", target))
			writer, err := os.Create(target)

			if err != nil {
				return err
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(target, os.FileMode(header.Mode))

			if err != nil {
				return err
			}

			writer.Close()
		default:
			bundleLog.Info(fmt.Sprintf("Unable to untar type : %c in file %s", header.Typeflag, target))
		}
	}
	return nil
}
