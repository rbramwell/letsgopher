package cmd

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type templateUninstallCmd struct {
	out     io.Writer
	name    string
	version string
	home    templ.Home
}

func newTemplateUninstallCmd(out io.Writer) *cobra.Command {
	remove := &templateUninstallCmd{out: out}

	cmd := &cobra.Command{
		Use:   "uninstall [NAME] [VERSION]",
		Short: "uninstall a template with a given name and version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "the template name", "the template version"); err != nil {
				return err
			}

			remove.name = args[0]
			remove.version = args[1]
			remove.home = templ.LetsGopherSettings.Home
			return remove.run()
		},
	}

	return cmd
}

func (r *templateUninstallCmd) run() error {
	err := deleteTemplateArchiveFile(r)
	if err != nil {
		return err
	}
	return removeTemplateLine(r.out, r.name, r.home)
}

func deleteTemplateArchiveFile(r *templateUninstallCmd) error {
	f, err := templ.LoadTemplatesFile(r.home.TemplatesFile())
	if err != nil {
		return err
	}

	template := f.Get(r.name, r.version)
	if template == nil {
		return fmt.Errorf("template with name %s and version %s hasn't been installed", r.name, r.version)
	}

	err = os.RemoveAll(template.ArchivePath)
	if err != nil {
		return fmt.Errorf("can't delete template archive %q", template.ArchivePath)
	}

	return nil
}

func removeTemplateLine(out io.Writer, name string, home templ.Home) error {
	templatesFile := home.TemplatesFile()
	r, err := templ.LoadTemplatesFile(templatesFile)
	if err != nil {
		return err
	}

	if !r.Remove(name) {
		return fmt.Errorf("no template named %q found", name)
	}
	if err := r.WriteFile(templatesFile, 0644); err != nil {
		return err
	}

	fmt.Fprintf(out, "%q has been removed from your templates\n", name)

	return nil
}
