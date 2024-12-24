package eui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/na2me/lib/resources"
)

func NewList(listOpts ...widget.ListOpt) *widget.List {
	return widget.NewList(
		append([]widget.ListOpt{
			widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(150, 0),
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionEnd,
					StretchVertical:    true,
					Padding:            widget.NewInsetsSimple(10),
				}),
			)),
			widget.ListOpts.ScrollContainerOpts(
				widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
					Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 100}),
					Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 100}),
					Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				}),
			),
			widget.ListOpts.HideHorizontalSlider(),
			widget.ListOpts.EntryFontFace(resources.Master.Fonts.BodyFace),
			widget.ListOpts.EntryColor(&widget.ListEntryColor{
				Selected:                   color.NRGBA{R: 0, G: 255, B: 0, A: 255},
				Unselected:                 color.NRGBA{R: 255, G: 255, B: 255, A: 255},
				SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255},
				SelectingBackground:        color.NRGBA{R: 130, G: 130, B: 130, A: 255},
				SelectingFocusedBackground: color.NRGBA{R: 130, G: 140, B: 170, A: 255},
				SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255},
				FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 100},
				DisabledUnselected:         color.NRGBA{R: 100, G: 100, B: 100, A: 255},
				DisabledSelected:           color.NRGBA{R: 100, G: 100, B: 100, A: 255},
				DisabledSelectedBackground: color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			}),
			widget.ListOpts.EntryLabelFunc(func(e interface{}) string { return "" }),
			widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(6)),
			widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
			widget.ListOpts.SliderOpts(
				widget.SliderOpts.Images(&widget.SliderTrackImage{
					Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
					Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				}, resources.Master.Button.Image),
				widget.SliderOpts.MinHandleSize(10),
				widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2)),
			),
			widget.ListOpts.AllowReselect(),
		}, listOpts...)...,
	)
}
