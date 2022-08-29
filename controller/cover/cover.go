package controller

const (
	ALT_FOR_COVER = "test"
)

// helper function for creating cover.xhtml file (/new-dir-###/EPUB/cover.xhtml)
func CoverXhtml(path, title string) string {
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
