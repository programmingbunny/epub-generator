package controller

import (
	"strconv"

	"github.com/programmingbunny/epub-generator/model"
)

func createChapterForTOC(allChapters []string) string {
	var test string
	for i := range allChapters {
		test = test + allChapters[i]
	}
	return test
}

// helper function for creating table of content (/new-dir-###/EPUB/bk01-toc.xhtml)
func CreateTOC(title, subtitle string, chapter model.Chapters) string {
	var stringTest []string
	for i := range chapter.Chapters {
		var singleTest string
		singleTest = `<li><a href="ch-` + strconv.Itoa(chapter.Chapters[i].ChapterNum) + `.xhtml">` + chapter.Chapters[i].Title + `</a></li>
					`
		stringTest = append(stringTest, singleTest)
	}
	testingThis := createChapterForTOC(stringTest)
	returnThis := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<!DOCTYPE html>
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="en"
		lang="en">
		<head>
			<title>` + title + `</title>
			<link rel="stylesheet" type="text/css" href="css/epub.css" />
		</head>
		<body>
			<h1>` + title + `</h1>
			<h1><i>` + subtitle + `</i></h1>
			<nav epub:type="toc" id="toc" role="doc-toc">
				<h2>Table of Contents</h2>
				<ol>
					` + testingThis + `
				</ol>
			</nav>
		</body>
	</html>`
	return returnThis
}
