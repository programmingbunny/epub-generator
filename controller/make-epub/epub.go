package controller

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	chapter "github.com/programmingbunny/epub-generator/controller/chapter"
	container "github.com/programmingbunny/epub-generator/controller/container"
	cover "github.com/programmingbunny/epub-generator/controller/cover"
	opf "github.com/programmingbunny/epub-generator/controller/package-opf"
	toc "github.com/programmingbunny/epub-generator/controller/toc"
	"github.com/programmingbunny/epub-generator/helpers"
	"github.com/programmingbunny/epub-generator/model"
)

const (
	// constants for directories
	NEW_DIRECTORY         = "new-dir-"
	EPUB                  = "/EPUB"
	META_INF              = "/META-INF"
	IMAGES                = "/images"
	COVERS                = "/covers"
	NO_FRONT_SLASH_COVERS = "covers/"
	MIMETYPE              = "/mimetype"
	COVER                 = "cover.xhtml"
	PACKAGE               = "package.opf"
	CONTAINER             = "container.xml"

	TOC            = "bk-toc"
	WRITE_MIMETYPE = "application/epub+zip"
	IMAGE_NAME     = "main-cover.jpg"
)

func CreateBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		bookId := params["bookId"]

		allChapters := chapter.GetChapters(bookId)
		newBook := chapter.GetBookDetails(bookId)

		fileName, err := creation(allChapters, newBook)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode("Something went wrong")
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode("Successfully created " + fileName)
	}
}

func creation(allChapters model.Chapters, newBook model.Book) (string, error) {
	name := helpers.NumberGen()

	cwd, _ := os.Getwd()

	// makes the parent directory (/new-dir-###/)
	err := makeNewDirectory(NEW_DIRECTORY + name)
	if err != nil {
		return "", err
	}

	// makes the EPUB directory within the parent directory (/new-dir-###/EPUB)
	makeNewDirectory(NEW_DIRECTORY + name + EPUB)
	if err != nil {
		return "", err
	}

	// makes the META-INF directory within the parent directory (/new-dir-###/META-INF)
	makeNewDirectory(NEW_DIRECTORY + name + META_INF)
	if err != nil {
		return "", err
	}

	// makes the EPUB/images directory within the parent directory (/new-dir-###/EPUB/images)
	makeNewDirectory(NEW_DIRECTORY + name + EPUB + IMAGES)
	if err != nil {
		return "", err
	}

	// makes the EPUB/covers directory (/new-dir-###/EPUB/covers)
	makeNewDirectory(NEW_DIRECTORY + name + EPUB + COVERS)
	if err != nil {
		return "", err
	}

	// create mimetype file in parent directory (/new-dir-###/mimetype)
	_, _, file, err := createFiles(cwd, NEW_DIRECTORY+name, "mimetype")
	if err != nil {
		return "", err
	}

	// opens & writes to mimetype file in parent directory (/new-dir-###/mimetype)
	openWriteFiles(file, NEW_DIRECTORY+name, MIMETYPE, WRITE_MIMETYPE)

	// create container.xml file in META-INF directory (/new-dir-###/META-INF/container.xml)
	_, _, file, err = createFiles(cwd, NEW_DIRECTORY+name+META_INF, CONTAINER)
	if err != nil {
		return "", err
	}

	// opens & writes to container.xml file in META-inf directory (/new-dir-###/META-INF/container.xml)
	openWriteFiles(file, NEW_DIRECTORY+name+META_INF, "/"+CONTAINER, container.ContainerXml())

	// adding cover image to EPUB/covers directory
	sourceFile, err := os.Open(newBook.BookCover)
	if err != nil {
		fmt.Println("Issue with os.open for newBook.BookCover: ", err)
	}
	defer sourceFile.Close()

	// create new cover image to EPUB/covers directory
	newFile, err := os.Create(NEW_DIRECTORY + name + EPUB + COVERS + "/" + IMAGE_NAME)
	if err != nil {
		fmt.Println("Issue with Create for EPUB/covers/IMAGE_NAME: ", err)
	}
	defer newFile.Close()

	// copy EPUB/covers image into directory
	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		log.Fatal(err)
	}

	// create cover.xhtml file in EPUB directory (/new-dir-###/EPUB/cover.xhtml)
	_, _, file, err = createFiles(cwd, NEW_DIRECTORY+name+EPUB, COVER)
	if err != nil {
		return "", err
	}

	// opens & writes to cover.xhtml file in META-inf directory (/new-dir-###/EPUB/cover.xhtml)
	openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/"+COVER, cover.CoverXhtml(NO_FRONT_SLASH_COVERS+IMAGE_NAME, newBook.Title))

	// create package.opf file in EPUB directory (/new-dir-###/EPUB/package.opf)
	_, _, file, err = createFiles(cwd, NEW_DIRECTORY+name+EPUB, PACKAGE)
	if err != nil {
		return "", err
	}

	// opens & writes to bk-toc.xhtml in EPUB directory (/new-dir-###/EPUB/bk-toc.xhtml
	openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/"+TOC+".xhtml", toc.CreateTOC(newBook.Title, newBook.Subtitle, allChapters))

	// opens & writes to package.opf file in EPUB directory (/new-dir-###/EPUB/package.opf)
	openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/"+PACKAGE, opf.EpubPackageOpf(NO_FRONT_SLASH_COVERS+IMAGE_NAME, newBook.Title, newBook.Author, allChapters))

	for i := range allChapters.Chapters {
		// adding chapter's header image
		if allChapters.Chapters[i].ImageLocation != "" {
			sourceFile, err := os.Open(allChapters.Chapters[i].ImageLocation)
			if err != nil {
				fmt.Println("Failing during Open of ImageLocation")
				log.Fatal(err)
			}
			defer sourceFile.Close()

			// create new chapter header to EPUB/images directory
			newFile, err := os.Create(NEW_DIRECTORY + name + EPUB + IMAGES + "/" + helpers.TrimImage(allChapters.Chapters[i].ImageLocation))
			if err != nil {
				log.Fatal(err)
			}
			defer newFile.Close()

			// copy EPUB/images image into directory
			_, err = io.Copy(newFile, sourceFile)
			if err != nil {
				log.Fatal(err)
			}
		}
		openWriteFiles(file, NEW_DIRECTORY+name+EPUB, "/ch-"+strconv.Itoa(allChapters.Chapters[i].ChapterNum)+".xhtml", chapter.CreateNewChapter(allChapters.Chapters[i]))
	}

	fmt.Println("Successfully created " + NEW_DIRECTORY + name)

	zipFolder(NEW_DIRECTORY+name+"/", newBook.Title+".epub")

	return NEW_DIRECTORY + name, nil
}

func zipFolder(source, target string) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
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
