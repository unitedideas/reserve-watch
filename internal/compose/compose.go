package compose

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"reserve-watch/internal/store"

	"github.com/fogleman/gg"
)

type Composer struct {
	templatesDir string
	outputDir    string
}

type ComposeInput struct {
	Topic      string
	SeriesName string
	Data       map[string]interface{}
}

type ComposeOutput struct {
	Blog       string
	LinkedIn   string
	Newsletter string
	Script     string
	ChartPNG   string
	OGPNG      string
}

func New(templatesDir, outputDir string) *Composer {
	os.MkdirAll(outputDir, 0755)
	return &Composer{
		templatesDir: templatesDir,
		outputDir:    outputDir,
	}
}

func (c *Composer) Compose(input ComposeInput, points []store.SeriesPoint) (*ComposeOutput, error) {
	if len(points) == 0 {
		return nil, fmt.Errorf("no data points provided")
	}

	latest := points[0]
	templateData := map[string]interface{}{
		"Title":             input.Data["title"],
		"SeriesName":        input.SeriesName,
		"CurrentValue":      fmt.Sprintf("%.2f", latest.Value),
		"CurrentDate":       latest.Date,
		"ChangeDescription": input.Data["change_description"],
		"Analysis":          input.Data["analysis"],
	}

	blogContent, err := c.renderTemplate("blog_note.tmpl", templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to render blog: %w", err)
	}

	linkedinContent, err := c.renderTemplate("linkedin.tmpl", templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to render linkedin: %w", err)
	}

	newsletterContent, err := c.renderTemplate("newsletter.tmpl", templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to render newsletter: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	chartPath := filepath.Join(c.outputDir, fmt.Sprintf("chart-%s.png", timestamp))

	if err := c.generateChart(points, input.SeriesName, chartPath); err != nil {
		return nil, fmt.Errorf("failed to generate chart: %w", err)
	}

	ogPath := filepath.Join(c.outputDir, fmt.Sprintf("og-%s.png", timestamp))
	if err := c.generateOGImage(latest.Value, input.SeriesName, latest.Date, ogPath); err != nil {
		return nil, fmt.Errorf("failed to generate OG image: %w", err)
	}

	return &ComposeOutput{
		Blog:       blogContent,
		LinkedIn:   linkedinContent,
		Newsletter: newsletterContent,
		Script:     fmt.Sprintf("The %s hit %.2f on %s. The reserve shift is accelerating.", input.SeriesName, latest.Value, latest.Date),
		ChartPNG:   chartPath,
		OGPNG:      ogPath,
	}, nil
}

func (c *Composer) renderTemplate(templateName string, data map[string]interface{}) (string, error) {
	templatePath := filepath.Join(c.templatesDir, templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *Composer) generateChart(points []store.SeriesPoint, seriesName, outputPath string) error {
	const width = 800
	const height = 600
	const margin = 60

	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawRectangle(margin, margin, float64(width-2*margin), float64(height-2*margin))
	dc.Stroke()

	if len(points) < 2 {
		dc.DrawStringAnchored("Insufficient data", float64(width/2), float64(height/2), 0.5, 0.5)
		return dc.SavePNG(outputPath)
	}

	minVal, maxVal := points[0].Value, points[0].Value
	for _, p := range points {
		if p.Value < minVal {
			minVal = p.Value
		}
		if p.Value > maxVal {
			maxVal = p.Value
		}
	}

	valueRange := maxVal - minVal
	if valueRange == 0 {
		valueRange = 1
	}

	chartWidth := float64(width - 2*margin)
	chartHeight := float64(height - 2*margin)

	dc.SetRGB(0.2, 0.4, 0.8)
	dc.SetLineWidth(2)

	for i := len(points) - 1; i >= 0; i-- {
		x := margin + chartWidth*(float64(len(points)-1-i)/float64(len(points)-1))
		y := margin + chartHeight*(1-(points[i].Value-minVal)/valueRange)

		if i == len(points)-1 {
			dc.MoveTo(x, y)
		} else {
			dc.LineTo(x, y)
		}
	}
	dc.Stroke()

	dc.SetRGB(0.1, 0.1, 0.1)
	if err := dc.LoadFontFace("C:\\Windows\\Fonts\\arial.ttf", 16); err != nil {
		dc.DrawStringAnchored(seriesName, float64(width/2), 30, 0.5, 0.5)
	} else {
		dc.DrawStringAnchored(seriesName, float64(width/2), 30, 0.5, 0.5)
	}

	dc.DrawStringAnchored(fmt.Sprintf("%.2f", minVal), margin-10, float64(height-margin), 1, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%.2f", maxVal), margin-10, float64(margin), 1, 0.5)

	return dc.SavePNG(outputPath)
}

func (c *Composer) generateOGImage(value float64, seriesName, date, outputPath string) error {
	const width = 1200
	const height = 630

	dc := gg.NewContext(width, height)

	dc.SetRGB(0.95, 0.95, 0.97)
	dc.Clear()

	dc.SetRGB(0.2, 0.3, 0.5)
	if err := dc.LoadFontFace("C:\\Windows\\Fonts\\arialbd.ttf", 48); err != nil {
		dc.DrawStringAnchored(seriesName, width/2, height/2-60, 0.5, 0.5)
	} else {
		dc.DrawStringAnchored(seriesName, width/2, height/2-60, 0.5, 0.5)
	}

	dc.SetRGB(0.1, 0.1, 0.1)
	if err := dc.LoadFontFace("C:\\Windows\\Fonts\\arial.ttf", 72); err != nil {
		dc.DrawStringAnchored(fmt.Sprintf("%.2f", value), width/2, height/2+20, 0.5, 0.5)
	} else {
		dc.DrawStringAnchored(fmt.Sprintf("%.2f", value), width/2, height/2+20, 0.5, 0.5)
	}

	dc.SetRGB(0.4, 0.4, 0.4)
	if err := dc.LoadFontFace("C:\\Windows\\Fonts\\arial.ttf", 24); err != nil {
		dc.DrawStringAnchored(date, width/2, height/2+80, 0.5, 0.5)
	} else {
		dc.DrawStringAnchored(date, width/2, height/2+80, 0.5, 0.5)
	}

	dc.SetRGB(0.3, 0.5, 0.8)
	dc.DrawStringAnchored("Reserve Watch", width/2, height-40, 0.5, 0.5)

	return dc.SavePNG(outputPath)
}
