package adrtest

import "testing"

func TestInitAdrRepository(t *testing.T) {
	w := newWorkspace(t)
	w.adr("init")
	w.ls("doc/adr")
	w.cat("doc/adr/0001-record-architecture-decisions.md")
	w.check()
}

func TestCreateFirstRecord(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "The First Decision")
	w.cat("doc/adr/0001-the-first-decision.md")
	w.check()
}

func TestCreateMultipleRecords(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Decision")
	w.adr("new", "Second Decision")
	w.adr("new", "Third Decision")
	w.ls("doc/adr")
	w.head(7, "doc/adr/0002-second-decision.md")
	w.head(7, "doc/adr/0003-third-decision.md")
	w.check()
}

func TestListRecords(t *testing.T) {
	w := newWorkspace(t)
	w.adr("list")
	w.adr("new", "first")
	w.adr("list")
	w.adr("new", "second")
	w.adr("list")
	w.adr("new", "third")
	w.adr("list")
	w.check()
}

func TestMustProvideATitleWhenCreatingNewAdr(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new")
	w.check()
}

func TestAvoidOctalNumbers(t *testing.T) {
	w := newWorkspace(t)
	for _, title := range []string{
		"First Decision", "Second Decision", "Third Decision", "Fourth Decision",
		"Fifth Decision", "Sixth Decision", "Seventh Decision", "Eighth Decision",
		"Ninth Decision",
	} {
		w.adr("new", title)
	}
	w.ls("doc/adr")
	w.head(7, "doc/adr/0009-ninth-decision.md")
	w.check()
}

func TestFunnyCharacters(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "Something About Node.JS")
	w.adr("new", "Slash/Slash/Slash/")
	w.adr("new", "--", "-Bar-")
	w.ls("doc/adr")
	w.check()
}

func TestEditAdrOnCreation(t *testing.T) {
	w := newWorkspace(t)
	w.useEditor("EDITOR", "EDITOR")
	w.setenv("VISUAL", "")
	w.adr("new", "created with EDITOR environment variable set")

	w.setenv("EDITOR", "")
	w.useEditor("VISUAL", "VISUAL")
	w.adr("new", "created with VISUAL environment variable set")

	w.useEditor("EDITOR", "EDITOR")
	w.useEditor("VISUAL", "VISUAL")
	w.adr("new", "uses setting of VISUAL if both VISUAL and EDITOR are set")
	w.check()
}

func TestAlternativeAdrDirectory(t *testing.T) {
	w := newWorkspace(t)
	w.adr("init", "alternative-dir")
	w.adr("new", "Example ADR")
	w.ls(".")
	w.ls("alternative-dir")
	w.adr("list")
	w.check()
}

func TestSearchForAdrDir(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Record")
	w.mkdir("subdir")
	w.cd("subdir")
	w.adr("new", "Second Record")
	w.adr("list")
	w.cd("..")
	w.adr("list")
	w.check()
}

func TestSearchForCustomAdrDir(t *testing.T) {
	w := newWorkspace(t)
	w.adr("init", "architecture-log")
	w.adr("new", "First Record")
	w.mkdir("subdir")
	w.cd("subdir")
	w.adr("new", "Second Record")
	w.adr("list")
	w.cd("..")
	w.adr("list")
	w.check()
}

func TestProjectSpecificTemplate(t *testing.T) {
	w := newWorkspace(t)
	w.adr("init", "adrs")
	w.mkdir("adrs/templates")
	w.write("adrs/templates/template.md", `# TITLE

Project specific template!

# Status

STATUS

# Info

ADR Number: NUMBER

Date: DATE
`)
	w.adr("new", "Aaa Bbb")
	w.adr("new", "Ccc Ddd")
	w.cat("adrs/0002-aaa-bbb.md")
	w.cat("adrs/0003-ccc-ddd.md")
	w.check()
}

func TestLinking(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Record")
	w.adr("new", "Second Record")
	w.adr("new", "Third Record")
	w.adr("link", "3", "Amends", "1", "Amended by")
	w.adr("link", "3", "Clarifies", "2", "Clarified by")
	w.head(12, "doc/adr/0001-first-record.md")
	w.head(12, "doc/adr/0002-second-record.md")
	w.head(14, "doc/adr/0003-third-record.md")
	w.check()
}

func TestLinkingNewRecords(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Record")
	w.adr("new", "Second Record")
	w.adr("new", "-l", "1:Amends:Amended by", "-l", "2:Clarifies:Clarified by", "Third Record")
	w.head(12, "doc/adr/0001-first-record.md")
	w.head(12, "doc/adr/0002-second-record.md")
	w.head(14, "doc/adr/0003-third-record.md")
	w.check()
}

func TestSupercedeExistingAdr(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Record")
	w.adr("new", "-s", "1", "Second Record")
	w.head(10, "doc/adr/0001-first-record.md")
	w.head(12, "doc/adr/0002-second-record.md")
	w.check()
}

func TestSupercedeMultipleAdrs(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Record")
	w.adr("new", "Second Record")
	w.adr("new", "-s", "1", "-s", "2", "Third Record")
	w.head(8, "doc/adr/0001-first-record.md")
	w.head(8, "doc/adr/0002-second-record.md")
	w.head(12, "doc/adr/0003-third-record.md")
	w.check()
}

func TestGenerateContents(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Decision")
	w.adr("new", "Second Decision")
	w.adr("new", "Third Decision")
	w.adr("generate", "toc")
	w.check()
}

func TestGenerateContentsWithPrefix(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Decision")
	w.adr("new", "Second Decision")
	w.adr("new", "Third Decision")
	w.adr("generate", "toc", "-p", "foo/doc/adr/")
	w.check()
}

func TestGenerateContentsWithHeaderAndFooter(t *testing.T) {
	w := newWorkspace(t)
	w.adr("new", "First Decision")
	w.adr("new", "Second Decision")
	w.adr("new", "Third Decision")
	w.write("intro.md", "An intro.\n\nMultiple paragraphs.\n")
	w.write("outro.md", "An outro.\n")
	w.adr("generate", "toc", "-i", "intro.md", "-o", "outro.md")
	w.check()
}

func TestGenerateGraph(t *testing.T) {
	w := newWorkspace(t)
	w.adr("init")
	w.adr("new", "An idea that seems good at the time")
	w.adr("new", "-s", "2", "A better idea")
	w.adr("new", "This will work")
	w.adr("new", "-s", "3", "The end")
	w.adr("generate", "graph")
	w.adr("generate", "graph", "-p", "http://example.com/", "-e", ".xxx")
	w.check()
}

func TestMigrateDateFormat(t *testing.T) {
	w := newWorkspace(t)
	w.adr("init")
	w.setenv("ADR_DATE", "12/01/1992")
	w.adr("new", "With Old Date Format")
	w.setenv("ADR_DATE", "1992-01-12")
	w.adr("new", "With Current Date Format")
	w.grep("Date:", "doc/adr/*")
	w.adr("upgrade-repository")
	w.grep("Date:", "doc/adr/*")
	w.check()
}

func TestAutocomplete(t *testing.T) {
	w := newWorkspace(t)
	w.adr("autocomplete", "adr")
	w.adr("autocomplete", "adr", "li")
	w.adr("autocomplete", "adr", "generate")
	w.adr("autocomplete", "adr", "generate", "to")
	w.adr("autocomplete", "adr", "ierhir")
	w.check()
}

func TestHelp(t *testing.T) {
	w := newWorkspace(t)
	w.adr("help")
	w.adr("help", "generate")
	w.adr("help", "generate", "toc")
	w.adr("nonsense")
	w.check()
}
