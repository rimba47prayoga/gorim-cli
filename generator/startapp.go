package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

var appFiles = map[string]string{
	"apps/%s/views.go": "templates/apps/views.go.tmpl",
	"apps/%s/serializer.go": "templates/apps/serializer.go.tmpl",
	"apps/%s/router.go": "templates/apps/router.go.tmpl",
}

func StartApp(appName string) {

	context := TemplateContext{AppName: appName}

	for path, tmplPath := range appFiles {
		path = fmt.Sprintf(path, appName)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Error creating directory:", err)
			continue
		}

		if err := CreateFileFromTemplate(path, tmplPath, context); err != nil {
			fmt.Println("Error creating file:", err)
		}
	}

	fmt.Println("App created successfully!")
}
