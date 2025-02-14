package decredmaterial

import (
	"image/color"

	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Modal struct {
	overlayColor color.NRGBA
	background   color.NRGBA
	list         *widget.List
	button       *widget.Clickable
	card         Card
	scroll       ListStyle
	padding      unit.Value

	isFloatTitle  bool
	isDisabled    bool
	showScrollBar bool
}

func (t *Theme) ModalFloatTitle() *Modal {
	mod := t.Modal()
	mod.isFloatTitle = true
	return mod
}

func (t *Theme) Modal() *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	m := &Modal{
		overlayColor: overlayColor,
		background:   t.Color.Surface,
		list: &widget.List{
			List: layout.List{Axis: layout.Vertical, Alignment: layout.Middle},
		},
		button:  new(widget.Clickable),
		card:    t.Card(),
		padding: unit.Dp(24),
	}

	m.scroll = t.List(m.list)

	return m
}

// Layout renders the modal widget to screen. The modal assumes the size of
// its content plus padding.
func (m *Modal) Layout(gtx layout.Context, widgets []layout.Widget) layout.Dimensions {
	mGtx := gtx
	if m.isDisabled {
		mGtx = gtx.Disabled()
	}
	dims := layout.Stack{Alignment: layout.Center}.Layout(mGtx,
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			fillMax(gtx, m.overlayColor, CornerRadius{})

			return m.button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				semantic.Button.Add(gtx.Ops)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			})
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			var widgetFuncs []layout.Widget
			var title layout.Widget

			if m.isFloatTitle && len(widgets) > 0 {
				title = widgets[0]
				widgetFuncs = append(widgetFuncs, widgets[1:]...)
			} else {
				widgetFuncs = append(widgetFuncs, widgets...)
			}

			gtx.Constraints.Max.X = gtx.Px(unit.Dp(380))
			inset := layout.Inset{
				Top:    unit.Dp(50),
				Bottom: unit.Dp(50),
			}
			return inset.Layout(gtx, func(gtx C) D {
				return LinearLayout{
					Orientation: layout.Vertical,
					Width:       WrapContent,
					Height:      WrapContent,
					Padding:     layout.UniformInset(m.padding),
					Alignment:   layout.Middle,
					Border: Border{
						Radius: Radius(14),
					},
					Direction:  layout.Center,
					Background: m.background,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if m.isFloatTitle && len(widgets) > 0 {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							if m.padding == unit.Dp(0) {
								return layout.UniformInset(m.padding).Layout(gtx, title)
							}

							inset := layout.Inset{
								Top:    unit.Dp(10),
								Bottom: unit.Dp(10),
							}
							return inset.Layout(gtx, title)
						}
						return D{}
					}),
					layout.Rigid(func(gtx C) D {
						mTB := unit.Dp(10)
						mLR := unit.Dp(0)
						if m.padding == unit.Dp(0) {
							mLR = mTB
						}
						inset := layout.Inset{
							Top:    mTB,
							Bottom: mTB,
							Left:   mLR,
							Right:  mLR,
						}
						if m.showScrollBar {
							return m.scroll.Layout(gtx, len(widgetFuncs), func(gtx C, i int) D {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return inset.Layout(gtx, widgetFuncs[i])
							})
						}
						list := &layout.List{Axis: layout.Vertical}
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return list.Layout(gtx, len(widgetFuncs), func(gtx C, i int) D {
							return inset.Layout(gtx, widgetFuncs[i])
						})
					}),
				)
			})
		}),
	)

	return dims
}

func (m *Modal) BackdropClicked(minimizable bool) bool {
	if minimizable {
		return m.button.Clicked()
	}

	return false
}

func (m *Modal) SetPadding(padding unit.Value) {
	m.padding = padding
}

func (m *Modal) ShowScrollbar(showScrollBar bool) {
	m.showScrollBar = showScrollBar
}

func (m *Modal) SetDisabled(disabled bool) {
	m.isDisabled = disabled
}
