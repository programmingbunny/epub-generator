package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	TITLE_NAME     = "Awesome Titles of Awesome"
	AUTHOR_NAME    = "Author McAwesome"
	PUBLISHER_NAME = "Publications, L.C."
	COPYRIGHT      = "Copyright © 2022 Publications"
	ALT_FOR_COVER  = "MARKED the Chronicles of a Fantasy Epic Series"

	// constants for directories
	NEW_DIRECTORY = "new-dir-"
	EPUB          = "/EPUB"
	META_INF      = "/META-INF"
	COVERS        = "/covers"
)

func main() {
	name := numberGen()

	cwd, _ := os.Getwd()

	// makes the parent directory
	err := makeNewDirectory(NEW_DIRECTORY + name)
	if err != nil {
		return
	}

	// makes the EPUB directory within the parent directory
	makeNewDirectory(NEW_DIRECTORY + name + EPUB)
	if err != nil {
		return
	}

	// makes the META-INF directory within the parent directory
	makeNewDirectory(NEW_DIRECTORY + name + META_INF)
	if err != nil {
		return
	}

	// makes the EPUB/covers directory within the EPUBS directory
	makeNewDirectory(NEW_DIRECTORY + name + EPUB + COVERS)
	if err != nil {
		return
	}

	// create mimetype file in parent directory
	path := filepath.Join(cwd, "new-dir-"+name, "mimetype")
	newFilePath := filepath.FromSlash(path)
	file, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// open mimetype file in parent directory to write to it
	file, err = os.OpenFile("new-dir-"+name+"/mimetype", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println("Could not open new-dir-" + name + " mimetype")
		return
	}

	defer file.Close()

	// write mimetype
	_, err2 := file.WriteString("application/epub+zip")
	if err2 != nil {
		fmt.Println("Could not write text to new-dir-" + name + " mimetype")
	}

	// create container.xml file in META-INF directory
	path = filepath.Join(cwd, "new-dir-"+name+META_INF, "container.xml")
	newFilePath = filepath.FromSlash(path)
	file, err = os.Create(newFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// open container.xml file in META-INF directory to write to it
	file, err = os.OpenFile("new-dir-"+name+META_INF+"/container.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println("Could not open new-dir-" + name + META_INF + "/container.xml")
		return
	}

	defer file.Close()

	// write container.xml
	_, err2 = file.WriteString(containerXml())
	if err2 != nil {
		fmt.Println("Could not write text to new-dir-" + name + META_INF + "/container.xml")
	}

	// adding cover image to EPUB/covers directory
	sourceFile, err := os.Open("./cover-test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()

	// create new cover image to EPUB/covers directory
	newFile, err := os.Create("new-dir-" + name + EPUB + COVERS + "/cover-test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// copy EPUB/covers image into directory
	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		log.Fatal(err)
	}

	// create cover.xhtml file in EPUB directory
	path = filepath.Join(cwd, "new-dir-"+name+EPUB, "cover.xhtml")
	newFilePath = filepath.FromSlash(path)
	file, err = os.Create(newFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// open cover.xhtml file in EPUB directory to write to it
	file, err = os.OpenFile("new-dir-"+name+EPUB+"cover.xhtml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println("Could not open new-dir-" + name + EPUB + "/cover.xhtml")
		return
	}

	defer file.Close()

	// write to cover.xhtml
	_, err2 = file.WriteString(coverXhtml(COVERS + "/cover-test.jpg"))
	if err2 != nil {
		fmt.Println("Could not write text to new-dir-" + name + EPUB + "/cover.xhtml")
	}

	// create package.opf file in EPUB directory
	path = filepath.Join(cwd, NEW_DIRECTORY+name+EPUB, "package.opf")
	newFilePath = filepath.FromSlash(path)
	file, err = os.Create(newFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// open package.opf file in EPUB directory to write to it
	file, err = os.OpenFile(NEW_DIRECTORY+name+EPUB+"/package.opf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println("Could not open new-dir-" + name + EPUB + "package.opf")
		return
	}

	defer file.Close()

	// write to package.opf
	_, err2 = file.WriteString(epubPackageOpf(COVERS + "/cover-test.jpg"))
	if err2 != nil {
		fmt.Println("Could not write text to " + NEW_DIRECTORY + name + EPUB + "/package.opf")
	}

}

func makeNewDirectory(path string) (err error) {
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		fmt.Println("Following error when trying to create"+path+"directory: ", err)
		return err
	}
	return nil
}

func numberGen() string {
	p, _ := rand.Prime(rand.Reader, 64)
	return p.String()
}

func coverXhtml(path string) string {
	returnThis := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="en" lang="en">
		<head>
			<title>` + TITLE_NAME + `</title>
			<style type="text/css">
				img{
					max-width:100%;
				}
			</style>
		</head>
		<body>
			<figure id="cover-image">
				<img role="doc-cover" src="` + path + `" alt="` + ALT_FOR_COVER + `" />
			</figure>
		</body>
	</html>`
	return returnThis
}

func epubPackageOpf(path string) string {
	returnThis := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
	<package xmlns="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/"
		xmlns:dcterms="http://purl.org/dc/terms/" version="3.0" xml:lang="en"
		unique-identifier="pub-identifier">
		<metadata>
			<dc:identifier id="pub-identifier">urn:isbn:123</dc:identifier>
			<dc:title id="pub-title">` + TITLE_NAME + `</dc:title>
			<dc:language id="pub-language">en</dc:language>
			<dc:date>2022-08-15</dc:date>
			<meta property="dcterms:modified">2012-10-24T15:30:00Z</meta>
			<dc:creator id="pub-creator12">` + AUTHOR_NAME + `</dc:creator>
			<dc:contributor>Fiona</dc:contributor>
			<dc:publisher>` + PUBLISHER_NAME + `</dc:publisher>
			<dc:rights>` + COPYRIGHT + `</dc:rights>
			<meta property="schema:accessMode">textual</meta>
			<meta property="schema:accessMode">visual</meta>
			<meta property="schema:accessModeSufficient">textual,visual</meta>
			<meta property="schema:accessModeSufficient">textual</meta>
			<meta property="schema:accessibilityHazard">none</meta>
			<meta property="schema:accessibilityFeature">tableOfContents</meta>
			<meta property="schema:accessibilityFeature">readingOrder</meta>
			<meta property="schema:accessibilityFeature">alternativeText</meta>
			<meta property="schema:accessibilitySummary">This EPUB Publication meets the requirements of the EPUB Accessibility specification with conformance to WCAG 2.0 Level AA. The publication is screen reader friendly.</meta>
			<link rel="dcterms:conformsTo" href="http://www.idpf.org/epub/a11y/accessibility-20170105.html#wcag-aa"/>
		</metadata>
		<manifest>
			<item id="cover" href="cover.xhtml" media-type="application/xhtml+xml"/>
			<item id="cover-image" properties="cover-image" href="` + path + `" media-type="image/jpeg"/>
		</manifest>
		<spine>
			<itemref idref="cover" linear="no"/>
			<itemref idref="htmltoc" linear="yes"/>
		</spine>
	</package>`
	return returnThis
}

func containerXml() string {
	returnThis := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
	<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container" version="1.0">
		<rootfiles>
			<rootfile full-path="EPUB/package.opf" media-type="application/oebps-package+xml"/>
		</rootfiles>
	</container>`
	return returnThis
}
