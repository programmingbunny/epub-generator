package controller

// helper function for creating container.xml (/new-dir-###/META-INF/container.xml)
func ContainerXml() string {
	returnThis := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
	<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container" version="1.0">
		<rootfiles>
			<rootfile full-path="EPUB/package.opf" media-type="application/oebps-package+xml"/>
		</rootfiles>
	</container>`
	return returnThis
}
