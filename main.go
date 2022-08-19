package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/programmingbunny/epub-generator/model"
)

const (
	PUBLISHER_NAME = "Publications, L.C."
	COPYRIGHT      = "Copyright © 2022 Publications"
	ALT_FOR_COVER  = "MARKED, the Chronicles of a Fantasy Epic Series"

	// constants for directories
	NEW_DIRECTORY         = "new-dir-"
	EPUB                  = "/EPUB"
	META_INF              = "/META-INF"
	COVERS                = "/covers"
	NO_FRONT_SLASH_COVERS = "covers/"
	MIMETYPE              = "/mimetype"

	WRITE_MIMETYPE = "application/epub+zip"
	IMAGE_NAME     = "cover-test.jpg"
)

func main() {
	id := "insertIdHereForTesting"
	response, err := http.Get("http://localhost:3000/book/" + id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var newBook model.Book
	json.Unmarshal(responseData, &newBook)

	name := numberGen()

	cwd, _ := os.Getwd()

	// makes the parent directory (/new-dir-###/)
	err = makeNewDirectory(NEW_DIRECTORY + name)
	if err != nil {
		return
	}

	// makes the EPUB directory within the parent directory (/new-dir-###/EPUB)
	makeNewDirectory(NEW_DIRECTORY + name + EPUB)
	if err != nil {
		return
	}

	// makes the META-INF directory within the parent directory (/new-dir-###/META-INF)
	makeNewDirectory(NEW_DIRECTORY + name + META_INF)
	if err != nil {
		return
	}

	// makes the EPUB/covers directory (/new-dir-###/EPUB/covers)
	makeNewDirectory(NEW_DIRECTORY + name + EPUB + COVERS)
	if err != nil {
		return
	}

	// create mimetype file in parent directory (/new-dir-###/mimetype)
	newFilePath, _, file, err := createFiles(cwd, NEW_DIRECTORY+name, "mimetype")
	if err != nil {
		return
	}

	// opens & writes to mimetype file in parent directory (/new-dir-###/mimetype)
	openWriteFiles(file, NEW_DIRECTORY+name, MIMETYPE, WRITE_MIMETYPE)

	// create container.xml file in META-INF directory (/new-dir-###/META-INF/container.xml)
	newFilePath, _, file, err = createFiles(cwd, NEW_DIRECTORY+name+META_INF, "container.xml")
	if err != nil {
		return
	}

	// opens & writes to container.xml file in META-inf directory (/new-dir-###/META-INF/container.xml)
	openWriteFiles(file, NEW_DIRECTORY+name+META_INF, "/container.xml", containerXml())

	// adding cover image to EPUB/covers directory
	sourceFile, err := os.Open("./cover-test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()

	// create new cover image to EPUB/covers directory
	newFile, err := os.Create(NEW_DIRECTORY + name + EPUB + COVERS + "/cover-test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// copy EPUB/covers image into directory
	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		log.Fatal(err)
	}

	// create cover.xhtml file in EPUB directory (/new-dir-###/EPUB/cover.xhtml)
	newFilePath, _, file, err = createFiles(cwd, NEW_DIRECTORY+name+EPUB, "cover.xhtml")
	if err != nil {
		return
	}

	// opens & writes to cover.xhtml file in META-inf directory (/new-dir-###/EPUB/cover.xhtml)
	openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/cover.xhtml", coverXhtml(NO_FRONT_SLASH_COVERS+IMAGE_NAME, newBook.Title))

	// create cover.xhtml file in EPUB directory (/new-dir-###/EPUB/package.opf)
	newFilePath, _, file, err = createFiles(cwd, NEW_DIRECTORY+name+EPUB, "package.opf")
	if err != nil {
		return
	}

	// opens & writes to package.opf file in EPUB directory (/new-dir-###/EPUB/package.opf)
	openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/package.opf", epubPackageOpf(NO_FRONT_SLASH_COVERS+IMAGE_NAME, newBook.Title, newBook.Author))

	fmt.Println("Successfully created ", newFilePath)
}

func makeNewDirectory(path string) (err error) {
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		fmt.Println("Following error when trying to create"+path+"directory: ", err)
		return err
	}
	return nil
}

// helper function for generating unique id for parent directory of epub file
func numberGen() string {
	p, _ := rand.Prime(rand.Reader, 64)
	return p.String()
}

// helper function for creating cover.xhtml file (/new-dir-###/EPUB/cover.xhtml)
func coverXhtml(path, title string) string {
	returnThis := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="en" lang="en">
		<head>
			<title>` + title + `</title>
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

// helper function for creating package.opf file (/new-dir-###/EPUB/package.opf)
func epubPackageOpf(path, title, author string) string {
	returnThis := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
	<package xmlns="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/"
		xmlns:dcterms="http://purl.org/dc/terms/" version="3.0" xml:lang="en"
		unique-identifier="pub-identifier">
		<metadata>
			<dc:identifier id="pub-identifier">urn:isbn:123</dc:identifier>
			<dc:title id="pub-title">` + title + `</dc:title>
			<dc:language id="pub-language">en</dc:language>
			<dc:date>2022-08-15</dc:date>
			<meta property="dcterms:modified">2012-10-24T15:30:00Z</meta>
			<dc:creator id="pub-creator12">` + author + `</dc:creator>
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

// helper function for creating container.xml (/new-dir-###/META-INF/container.xml)
func containerXml() string {
	returnThis := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
	<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container" version="1.0">
		<rootfiles>
			<rootfile full-path="EPUB/package.opf" media-type="application/oebps-package+xml"/>
		</rootfiles>
	</container>`
	return returnThis
}

func createFiles(cwd, directory, fileName string) (string, string, *os.File, error) {
	path := filepath.Join(cwd, directory, fileName)
	newFilePath := filepath.FromSlash(path)
	file, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println(err)
		return "", "", nil, err
	}
	defer file.Close()

	return newFilePath, path, file, nil
}

func openWriteFiles(file *os.File, path, fileName, write string) error {
	file, err := os.OpenFile(path+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println("Could not open ", path+fileName+": ", err)
		return nil
	}
	defer file.Close()

	_, err2 := file.WriteString(write)
	if err2 != nil {
		fmt.Println("Could not write text to "+path+fileName+": ", err)
		return nil
	}
	return nil
}
