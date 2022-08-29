package controller

import (
	"strconv"

	"github.com/programmingbunny/epub-generator/helpers"
	"github.com/programmingbunny/epub-generator/model"
)

const (
	PUBLISHER_NAME = "Publications, L.C."
	COPYRIGHT      = "Copyright © 2022 Publications"

	TOC = "bk-toc"
)

// generates html for each chapter of the book to be used in package.opf file creation
func CreateItemIdForPackageOpf(chapter model.Chapters) (string, string) {
	var itemId string
	var arrayOfItemIds []string
	for i := range chapter.Chapters {
		singleItemId := helpers.NumberGen()
		itemId = itemId + `<item id="id-` + singleItemId + `" href="ch-` + strconv.Itoa(chapter.Chapters[i].ChapterNum) + `.xhtml" media-type="application/xhtml+xml"/>
					`
		arrayOfItemIds = append(arrayOfItemIds, singleItemId)
	}

	var stringForItemIds string
	for i := range arrayOfItemIds {
		stringForItemIds = stringForItemIds + `<itemref idref="id-` + arrayOfItemIds[i] + `"/>
					`
	}
	return itemId, stringForItemIds
}

// helper function for creating package.opf file (/new-dir-###/EPUB/package.opf)
func EpubPackageOpf(path, title, author string, chapter model.Chapters) string {
	chapterList, chapterIds := CreateItemIdForPackageOpf(chapter)
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
			<item id="htmltoc" properties="nav" media-type="application/xhtml+xml" href="` + TOC + `.xhtml"/>
			<item id="cover" href="cover.xhtml" media-type="application/xhtml+xml"/>
			<item id="cover-image" properties="cover-image" href="` + path + `" media-type="image/jpeg"/>
			` + chapterList + `
		</manifest>
		<spine>
			<itemref idref="cover" linear="no"/>
			<itemref idref="htmltoc" linear="yes"/>
			` + chapterIds + `
		</spine>
	</package>`
	return returnThis
}
