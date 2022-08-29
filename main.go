package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/programmingbunny/epub-generator/controller/chapter"
	container "github.com/programmingbunny/epub-generator/controller/container"
	cover "github.com/programmingbunny/epub-generator/controller/cover"
	opf "github.com/programmingbunny/epub-generator/controller/package-opf"
	toc "github.com/programmingbunny/epub-generator/controller/toc"
	"github.com/programmingbunny/epub-generator/helpers"
)

const (
	// constants for directories
	NEW_DIRECTORY         = "new-dir-"
	EPUB                  = "/EPUB"
	META_INF              = "/META-INF"
	COVERS                = "/covers"
	NO_FRONT_SLASH_COVERS = "covers/"
	MIMETYPE              = "/mimetype"

	TOC            = "bk-toc"
	WRITE_MIMETYPE = "application/epub+zip"
	IMAGE_NAME     = "cover-test.jpg"
)

func main() {
	newBook := chapter.GetBookDetails("630c424a0b3339afae9fcbf0")

	allChapters := chapter.GetChapters("630c424a0b3339afae9fcbf0")

	name := helpers.NumberGen()

	cwd, _ := os.Getwd()

	// makes the parent directory (/new-dir-###/)
	err := makeNewDirectory(NEW_DIRECTORY + name)
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
	openWriteFiles(file, NEW_DIRECTORY+name+META_INF, "/container.xml", container.ContainerXml())

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
	openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/cover.xhtml", cover.CoverXhtml(NO_FRONT_SLASH_COVERS+IMAGE_NAME, newBook.Title))

	// create cover.xhtml file in EPUB directory (/new-dir-###/EPUB/package.opf)
	newFilePath, _, file, err = createFiles(cwd, NEW_DIRECTORY+name+EPUB, "package.opf")
	if err != nil {
		return
	}

	// opens & writes to bk-toc.xhtml in EPUB directory (/new-dir-###/EPUB/bk-toc.xhtml
	openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/"+TOC+".xhtml", toc.CreateTOC(newBook.Title, newBook.Subtitle, allChapters))

	// opens & writes to package.opf file in EPUB directory (/new-dir-###/EPUB/package.opf)
	openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/package.opf", opf.EpubPackageOpf(NO_FRONT_SLASH_COVERS+IMAGE_NAME, newBook.Title, newBook.Author, allChapters))

	for i := range allChapters.Chapters {
		openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/ch-"+strconv.Itoa(allChapters.Chapters[i].ChapterNum)+".xhtml", chapter.CreateNewChapter(allChapters.Chapters[i]))
	}

	fmt.Println("Successfully created ", newFilePath)

}

// create directory/files
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

func makeNewDirectory(path string) (err error) {
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		fmt.Println("Following error when trying to create"+path+"directory: ", err)
		return err
	}
	return nil
}
