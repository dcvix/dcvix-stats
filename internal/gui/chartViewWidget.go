//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

package gui

import (
	"bytes"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dcvix/dcvix-stats/internal/charts"
)

// Ensure ChartView implements fyne.Widget
// emit compile error if interface is not implemented (cast nil to interface)
var _ fyne.Widget = (*ChartView)(nil)

// const minSizeW = 512
// const minSizeH = 284
const minSizeW = 500
const minSizeH = 200

// Add more to implement more interfaces for example:
// var _ fyne.Draggable = (*ChartView)(nil)

// ChartView is a widget that displays an image.
type ChartView struct {
	widget.BaseWidget

	// locker sync.Mutex
	img        *canvas.Image
	metrics    []string
	values     [][]float64
	timeStamps []string
}

// NewChartView creates a new ChartView widget. It implements fyne.Widget.
func NewChartView(metrics []string, values [][]float64, timeStamps []string) *ChartView {
	c := &ChartView{
		img:        canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 1, 1))), // Needs a placeholder image
		metrics:    metrics,
		values:     values,
		timeStamps: timeStamps,
	}

	c.ExtendBaseWidget(c)
	return c
}

// Set a sane minimal size or the graph will be unreadable
func (w *ChartView) MinSize() fyne.Size {
	return fyne.NewSize(minSizeW, minSizeH)
}

func (c *ChartView) Resize(s fyne.Size) {
	// Make the Chart resolution bigger than widget size for a clearer graph, fyne canvas will resize it
	if needRerender(float32(c.img.Image.Bounds().Dx()), float32(c.img.Image.Bounds().Dy()), s.Height, s.Width) {
		imgSize := fyne.NewSize(s.Width*1.2, s.Height*1.2)
		c.img.Image = c.GenerateChart(imgSize)
		c.img.FillMode = canvas.ImageFillStretch
	}
	c.BaseWidget.Resize(s)
}

// GenerateChart generate a chart image used by this widget.
func (c *ChartView) GenerateChart(s fyne.Size) image.Image {
	Width, Height := s.Width, s.Height

	// Handle zero or very small dimensions
	if s.Width <= 1 || s.Height <= 1 {
		Width = 600
		Height = 300 // sensible default 2:1 ratio
	} else {
		// Maintain aspect ratio while enforcing minimum width
		if Width < 600 {
			ratio := float32(s.Width) / float32(s.Height)
			Width = 600
			Height = Width / ratio
		}
	}

	chartImageBuff := charts.Chart(c.metrics, c.values, c.timeStamps, Width, Height)
	chartImageReader := bytes.NewReader(chartImageBuff)
	chartImage, _, _ := image.Decode(chartImageReader)
	return chartImage
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (c *ChartView) CreateRenderer() fyne.WidgetRenderer {
	co := container.NewAdaptiveGrid(1, c.img)
	return widget.NewSimpleRenderer(co)
}

// Re-render graphs with new data
func (c *ChartView) RefreshData(values [][]float64, timeStamps []string) {
	c.values = values
	c.timeStamps = timeStamps
	c.img.Image = c.GenerateChart(c.Size())
}

// needRerender checks if chart image needs to be re-rendered for new widget size
func needRerender(dx, dy, wdgDx, wdgDy float32) bool {
	if wdgDx > dx || wdgDy > dy || wdgDx < (dx*0.5) || wdgDy < (dy*0.5) {
		// fmt.Printf("Rerender: dx:%v dy:%v wx:%v wy:%v\n", dx, dy, wdgDx, wdgDy)
		return true
	}
	return false
}
