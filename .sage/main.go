package main

import (
	"context"
	"os"
	"path/filepath"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/sgtool"
	"go.einride.tech/sage/tools/sgconvco"
	"go.einride.tech/sage/tools/sggit"
	"go.einride.tech/sage/tools/sggo"
	"go.einride.tech/sage/tools/sggolangcilint"
	"go.einride.tech/sage/tools/sggolicenses"
	"go.einride.tech/sage/tools/sggoreview"
	"go.einride.tech/sage/tools/sgmdformat"
	"go.einride.tech/sage/tools/sgyamlfmt"
)

func main() {
	sg.GenerateMakefiles(
		sg.Makefile{
			Path:          sg.FromGitRoot("Makefile"),
			DefaultTarget: Default,
		},
	)
}

func Default(ctx context.Context) error {
	sg.Deps(ctx, ConvcoCheck, FormatMarkdown, FormatYaml, GoGenerate, GenerateTestdata)
	sg.Deps(ctx, GoLint, GoReview)
	sg.Deps(ctx, GoTest)
	sg.Deps(ctx, GoModTidy)
	sg.Deps(ctx, GoLicenses, GitVerifyNoDiff)
	return nil
}

func GoModTidy(ctx context.Context) error {
	sg.Logger(ctx).Println("tidying Go module files...")
	return sg.Command(ctx, "go", "mod", "tidy", "-v").Run()
}

func GoTest(ctx context.Context) error {
	sg.Logger(ctx).Println("running Go tests...")
	return sggo.TestCommand(ctx).Run()
}

func GoReview(ctx context.Context) error {
	sg.Logger(ctx).Println("reviewing Go files...")
	return sggoreview.Run(ctx)
}

func GoLint(ctx context.Context) error {
	sg.Logger(ctx).Println("linting Go files...")
	return sggolangcilint.Run(ctx)
}

func GoLicenses(ctx context.Context) error {
	sg.Logger(ctx).Println("checking Go licenses...")
	return sggolicenses.Check(ctx)
}

func FormatMarkdown(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting Markdown files...")
	return sgmdformat.Command(ctx).Run()
}

func FormatYaml(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting Yaml files...")
	return sgyamlfmt.Run(ctx)
}

func ConvcoCheck(ctx context.Context) error {
	sg.Logger(ctx).Println("checking git commits...")
	return sgconvco.Command(ctx, "check", "origin/master..HEAD").Run()
}

func GitVerifyNoDiff(ctx context.Context) error {
	sg.Logger(ctx).Println("verifying that git has no diff...")
	return sggit.VerifyNoDiff(ctx)
}

func GoGenerate(ctx context.Context) error {
	sg.Deps(ctx, Mockgen, Stringer)
	sg.Logger(ctx).Println("generating Go code...")
	return sg.Command(ctx, "go", "generate", "./...").Run()
}

func Mockgen(ctx context.Context) error {
	sg.Logger(ctx).Println("installing mockgen...")
	_, err := sgtool.GoInstallWithModfile(ctx, "github.com/golang/mock/mockgen", sg.FromGitRoot("go.mod"))
	return err
}

func Stringer(ctx context.Context) error {
	sg.Logger(ctx).Println("installing stringer...")
	_, err := sgtool.GoInstallWithModfile(ctx, "golang.org/x/tools/cmd/stringer", sg.FromGitRoot("go.mod"))
	return err
}

func GenerateTestdata(ctx context.Context) error {
	sg.Logger(ctx).Println("generating testdata...")
	// don't use "sg.FromGitRoot" in paths to avoid embedding user paths in generated files
	cmd := sg.Command(
		ctx,
		"go",
		"run",
		"cmd/cantool/main.go",
		"generate",
		"testdata/dbc",
		"testdata/gen/go",
	)
	cmd.Dir = sg.FromGitRoot()
	return cmd.Run()
}

func BuildIntegrationTests(ctx context.Context) error {
	sg.Logger(ctx).Println("building integration test...")
	testDir := sg.FromGitRoot("build", "tests")
	if err := os.MkdirAll(testDir, 0o775); err != nil {
		return err
	}
	return sg.Command(
		ctx,
		"go",
		"test",
		"-tags=integration",
		"-c",
		sg.FromGitRoot("pkg", "candevice"),
		"-o",
		filepath.Join(testDir, "candevice.test"),
	).Run()
}
