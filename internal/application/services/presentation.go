package services

import (
	"archive/zip"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PresentationGenerator handles converting Marp markdown to PowerPoint
type PresentationGenerator struct{}

// NewPresentationGenerator creates a new presentation generator
func NewPresentationGenerator() *PresentationGenerator {
	return &PresentationGenerator{}
}

// Slide represents a single slide with its content
type Slide struct {
	Title       string
	Content     []string
	Notes       string
	HasDiagram  bool
	DiagramType string
	DiagramText string
	IsCodeSlide bool
	CodeContent string
}

// GenerateFromMarp converts a Marp markdown file to PowerPoint
func (pg *PresentationGenerator) GenerateFromMarp(inputPath, outputPath string) error {
	// Parse the Marp markdown
	slides, err := pg.parseMarpMarkdown(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse markdown: %w", err)
	}

	// Create PPTX file
	if err := pg.createPPTX(slides, outputPath); err != nil {
		return fmt.Errorf("failed to create PPTX: %w", err)
	}

	return nil
}

// createPPTX creates a PPTX file from slides
func (pg *PresentationGenerator) createPPTX(slides []Slide, outputPath string) error {
	// Create a new zip file
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Write [Content_Types].xml
	if err := pg.writeContentTypes(zipWriter, len(slides)); err != nil {
		return err
	}

	// Write _rels/.rels
	if err := pg.writeRels(zipWriter); err != nil {
		return err
	}

	// Write ppt/_rels/presentation.xml.rels
	if err := pg.writePresentationRels(zipWriter, len(slides)); err != nil {
		return err
	}

	// Write ppt/presentation.xml
	if err := pg.writePresentation(zipWriter, len(slides)); err != nil {
		return err
	}

	// Write slides
	for i, slide := range slides {
		if err := pg.writeSlide(zipWriter, i+1, slide); err != nil {
			return err
		}
		if err := pg.writeSlideRels(zipWriter, i+1); err != nil {
			return err
		}
	}

	// Write slide layouts and masters (minimal)
	if err := pg.writeSlideLayouts(zipWriter); err != nil {
		return err
	}

	return nil
}

// writeContentTypes writes the [Content_Types].xml file
func (pg *PresentationGenerator) writeContentTypes(zw *zip.Writer, slideCount int) error {
	w, err := zw.Create("[Content_Types].xml")
	if err != nil {
		return err
	}

	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>`

	for i := 1; i <= slideCount; i++ {
		xml += fmt.Sprintf(`
  <Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, i)
	}

	xml += `
  <Override PartName="/ppt/slideLayouts/slideLayout1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
  <Override PartName="/ppt/slideMasters/slideMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>
</Types>`

	_, err = w.Write([]byte(xml))
	return err
}

// writeRels writes the _rels/.rels file
func (pg *PresentationGenerator) writeRels(zw *zip.Writer) error {
	w, err := zw.Create("_rels/.rels")
	if err != nil {
		return err
	}

	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
</Relationships>`

	_, err = w.Write([]byte(xml))
	return err
}

// writePresentationRels writes the ppt/_rels/presentation.xml.rels file
func (pg *PresentationGenerator) writePresentationRels(zw *zip.Writer, slideCount int) error {
	w, err := zw.Create("ppt/_rels/presentation.xml.rels")
	if err != nil {
		return err
	}

	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`

	for i := 1; i <= slideCount; i++ {
		xml += fmt.Sprintf(`
  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide%d.xml"/>`, i, i)
	}

	xml += fmt.Sprintf(`
  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster1.xml"/>
</Relationships>`, slideCount+1)

	_, err = w.Write([]byte(xml))
	return err
}

// writePresentation writes the ppt/presentation.xml file
func (pg *PresentationGenerator) writePresentation(zw *zip.Writer, slideCount int) error {
	w, err := zw.Create("ppt/presentation.xml")
	if err != nil {
		return err
	}

	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:sldMasterIdLst>
    <p:sldMasterId id="2147483648" r:id="rId` + fmt.Sprintf("%d", slideCount+1) + `"/>
  </p:sldMasterIdLst>
  <p:sldIdLst>`

	for i := 1; i <= slideCount; i++ {
		xml += fmt.Sprintf(`
    <p:sldId id="%d" r:id="rId%d"/>`, 255+i, i)
	}

	xml += `
  </p:sldIdLst>
  <p:sldSz cx="9144000" cy="6858000" type="screen4x3"/>
  <p:notesSz cx="6858000" cy="9144000"/>
</p:presentation>`

	_, err = w.Write([]byte(xml))
	return err
}

// writeSlide writes a slide XML file
func (pg *PresentationGenerator) writeSlide(zw *zip.Writer, slideNum int, slide Slide) error {
	w, err := zw.Create(fmt.Sprintf("ppt/slides/slide%d.xml", slideNum))
	if err != nil {
		return err
	}

	xmlContent := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr>
        <a:xfrm>
          <a:off x="0" y="0"/>
          <a:ext cx="0" cy="0"/>
          <a:chOff x="0" y="0"/>
          <a:chExt cx="0" cy="0"/>
        </a:xfrm>
      </p:grpSpPr>`

	shapeId := 2

	// Add title shape
	if slide.Title != "" {
		xmlContent += fmt.Sprintf(`
      <p:sp>
        <p:nvSpPr>
          <p:cNvPr id="%d" name="Title"/>
          <p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
          <p:nvPr><p:ph type="ctrTitle"/></p:nvPr>
        </p:nvSpPr>
        <p:spPr>
          <a:xfrm>
            <a:off x="914400" y="685800"/>
            <a:ext cx="7315200" cy="1143000"/>
          </a:xfrm>
        </p:spPr>
        <p:txBody>
          <a:bodyPr/>
          <a:lstStyle/>
          <a:p>
            <a:r>
              <a:rPr lang="en-US" sz="3200" b="1"/>
              <a:t>%s</a:t>
            </a:r>
          </a:p>
        </p:txBody>
      </p:sp>`, shapeId, escapeXML(slide.Title))
		shapeId++
	}

	// Add content shape
	if len(slide.Content) > 0 {
		xmlContent += fmt.Sprintf(`
      <p:sp>
        <p:nvSpPr>
          <p:cNvPr id="%d" name="Content"/>
          <p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
          <p:nvPr><p:ph type="body" idx="1"/></p:nvPr>
        </p:nvSpPr>
        <p:spPr>
          <a:xfrm>
            <a:off x="914400" y="1828800"/>
            <a:ext cx="7315200" cy="4572000"/>
          </a:xfrm>
        </p:spPr>
        <p:txBody>
          <a:bodyPr/>
          <a:lstStyle/>`, shapeId)

		for _, line := range slide.Content {
			level := 0
			if strings.HasPrefix(line, "  - ") {
				level = 1
				line = strings.TrimPrefix(line, "  - ")
			} else if strings.HasPrefix(line, "- ") {
				level = 0
				line = strings.TrimPrefix(line, "- ")
			}

			xmlContent += fmt.Sprintf(`
          <a:p>
            <a:pPr lvl="%d"/>
            <a:r>
              <a:rPr lang="en-US" sz="1800"/>
              <a:t>%s</a:t>
            </a:r>
          </a:p>`, level, escapeXML(line))
		}

		xmlContent += `
        </p:txBody>
      </p:sp>`
		shapeId++
	}

	// Add diagram placeholder
	if slide.HasDiagram {
		diagramText := fmt.Sprintf("[%s diagram]\n\n%s", slide.DiagramType, slide.DiagramText)
		if len(diagramText) > 500 {
			diagramText = diagramText[:500] + "..."
		}
		xmlContent += fmt.Sprintf(`
      <p:sp>
        <p:nvSpPr>
          <p:cNvPr id="%d" name="Diagram"/>
          <p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
          <p:nvPr/>
        </p:nvSpPr>
        <p:spPr>
          <a:xfrm>
            <a:off x="914400" y="1828800"/>
            <a:ext cx="7315200" cy="4114800"/>
          </a:xfrm>
          <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
          <a:solidFill><a:srgbClr val="F0F0F0"/></a:solidFill>
          <a:ln><a:solidFill><a:srgbClr val="000000"/></a:solidFill></a:ln>
        </p:spPr>
        <p:txBody>
          <a:bodyPr/>
          <a:lstStyle/>
          <a:p>
            <a:r>
              <a:rPr lang="en-US" sz="1000"/>
              <a:t>%s</a:t>
            </a:r>
          </a:p>
        </p:txBody>
      </p:sp>`, shapeId, escapeXML(diagramText))
		shapeId++
	}

	// Add code content
	if slide.IsCodeSlide && slide.CodeContent != "" {
		xmlContent += fmt.Sprintf(`
      <p:sp>
        <p:nvSpPr>
          <p:cNvPr id="%d" name="Code"/>
          <p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
          <p:nvPr/>
        </p:nvSpPr>
        <p:spPr>
          <a:xfrm>
            <a:off x="685800" y="1828800"/>
            <a:ext cx="7772400" cy="4572000"/>
          </a:xfrm>
          <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
          <a:solidFill><a:srgbClr val="282C34"/></a:solidFill>
          <a:ln><a:solidFill><a:srgbClr val="646464"/></a:solidFill></a:ln>
        </p:spPr>
        <p:txBody>
          <a:bodyPr/>
          <a:lstStyle/>
          <a:p>
            <a:r>
              <a:rPr lang="en-US" sz="1200"/>
              <a:t>%s</a:t>
            </a:r>
          </a:p>
        </p:txBody>
      </p:sp>`, shapeId, escapeXML(slide.CodeContent))
	}

	xmlContent += `
    </p:spTree>
  </p:cSld>
  <p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>
</p:sld>`

	_, err = w.Write([]byte(xmlContent))
	return err
}

// writeSlideRels writes slide relationship files
func (pg *PresentationGenerator) writeSlideRels(zw *zip.Writer, slideNum int) error {
	w, err := zw.Create(fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", slideNum))
	if err != nil {
		return err
	}

	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
</Relationships>`

	_, err = w.Write([]byte(xml))
	return err
}

// writeSlideLayouts writes minimal slide layout and master files
func (pg *PresentationGenerator) writeSlideLayouts(zw *zip.Writer) error {
	// Write slide layout
	w1, err := zw.Create("ppt/slideLayouts/slideLayout1.xml")
	if err != nil {
		return err
	}

	layoutXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldLayout xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" type="blank" preserve="1">
  <p:cSld name="Blank">
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr/>
    </p:spTree>
  </p:cSld>
  <p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>
</p:sldLayout>`

	if _, err := w1.Write([]byte(layoutXML)); err != nil {
		return err
	}

	// Write slide layout rels
	w2, err := zw.Create("ppt/slideLayouts/_rels/slideLayout1.xml.rels")
	if err != nil {
		return err
	}

	layoutRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="../slideMasters/slideMaster1.xml"/>
</Relationships>`

	if _, err := w2.Write([]byte(layoutRels)); err != nil {
		return err
	}

	// Write slide master
	w3, err := zw.Create("ppt/slideMasters/slideMaster1.xml")
	if err != nil {
		return err
	}

	masterXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr/>
    </p:spTree>
  </p:cSld>
  <p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>
  <p:sldLayoutIdLst>
    <p:sldLayoutId id="2147483649" r:id="rId1"/>
  </p:sldLayoutIdLst>
</p:sldMaster>`

	if _, err := w3.Write([]byte(masterXML)); err != nil {
		return err
	}

	// Write slide master rels
	w4, err := zw.Create("ppt/slideMasters/_rels/slideMaster1.xml.rels")
	if err != nil {
		return err
	}

	masterRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
</Relationships>`

	_, err = w4.Write([]byte(masterRels))
	return err
}

// parseMarpMarkdown parses a Marp markdown file into slides
func (pg *PresentationGenerator) parseMarpMarkdown(filePath string) ([]Slide, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var slides []Slide
	var currentSlide *Slide
	var inFrontmatter bool
	var inCodeBlock bool
	var inNotes bool
	var codeBlockLang string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Handle frontmatter
		if strings.HasPrefix(line, "---") {
			if currentSlide == nil && !inFrontmatter {
				inFrontmatter = true
				continue
			} else if inFrontmatter {
				inFrontmatter = false
				continue
			} else {
				// New slide separator
				if currentSlide != nil {
					slides = append(slides, *currentSlide)
				}
				currentSlide = &Slide{}
				inNotes = false
				continue
			}
		}

		if inFrontmatter {
			continue
		}

		if currentSlide == nil {
			continue
		}

		// Handle code blocks
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				inCodeBlock = true
				codeBlockLang = strings.TrimPrefix(line, "```")
				if codeBlockLang == "mermaid" {
					currentSlide.HasDiagram = true
					currentSlide.DiagramType = "mermaid"
				}
				continue
			} else {
				inCodeBlock = false
				codeBlockLang = ""
				continue
			}
		}

		if inCodeBlock {
			if currentSlide.HasDiagram && currentSlide.DiagramType == "mermaid" {
				currentSlide.DiagramText += line + "\n"
			} else {
				currentSlide.IsCodeSlide = true
				currentSlide.CodeContent += line + "\n"
			}
			continue
		}

		// Handle notes section
		if strings.HasPrefix(line, "Notes:") {
			inNotes = true
			continue
		}

		if inNotes {
			if strings.HasPrefix(line, "- ") {
				currentSlide.Notes += strings.TrimPrefix(line, "- ") + " "
			}
			continue
		}

		// Handle title (h1)
		if strings.HasPrefix(line, "# ") {
			currentSlide.Title = strings.TrimPrefix(line, "# ")
			continue
		}

		// Handle subtitle (h2)
		if strings.HasPrefix(line, "## ") {
			currentSlide.Content = append(currentSlide.Content, strings.TrimPrefix(line, "## "))
			continue
		}

		// Handle alt-text
		if strings.HasPrefix(line, "Alt-text:") {
			continue
		}

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Add bullet points and other content
		currentSlide.Content = append(currentSlide.Content, line)
	}

	// Add the last slide
	if currentSlide != nil {
		slides = append(slides, *currentSlide)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return slides, nil
}

// escapeXML escapes special XML characters
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
