package page

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const VerifyMessagePageID = "VerifyMessage"

type VerifyMessagePage struct {
	*load.Load

	addressEditor          decredmaterial.Editor
	messageEditor          decredmaterial.Editor
	signatureEditor        decredmaterial.Editor
	clearBtn, verifyButton decredmaterial.Button
	verifyMessage          decredmaterial.Label
	EnableEditorSwitch     bool

	verifyMessageStatus *decredmaterial.Icon

	backButton     decredmaterial.IconButton
	infoButton     decredmaterial.IconButton
	addressIsValid bool
	isEnabled      bool
}

func NewVerifyMessagePage(l *load.Load) *VerifyMessagePage {
	pg := &VerifyMessagePage{
		Load:               l,
		verifyMessage:      l.Theme.Body1(""),
		EnableEditorSwitch: false,
	}

	pg.addressEditor = l.Theme.Editor(new(widget.Editor), "Address")
	pg.addressEditor.Editor.SingleLine = true
	pg.addressEditor.Editor.Submit = true

	pg.messageEditor = l.Theme.Editor(new(widget.Editor), "Message")
	pg.messageEditor.Editor.SingleLine = true
	pg.messageEditor.Editor.Submit = true

	pg.signatureEditor = l.Theme.Editor(new(widget.Editor), "Signature")
	pg.signatureEditor.Editor.Submit = true

	pg.verifyButton = l.Theme.Button("Verify message")
	pg.verifyButton.Font.Weight = text.Medium

	pg.clearBtn = l.Theme.OutlineButton("Clear all")
	pg.clearBtn.Font.Weight = text.Medium

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *VerifyMessagePage) ID() string {
	return VerifyMessagePageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *VerifyMessagePage) OnNavigatedTo() {
	pg.addressEditor.Editor.Focus()
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *VerifyMessagePage) Layout(gtx C) D {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Verify message",
			BackButton: pg.backButton,
			InfoButton: pg.infoButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(pg.description()),
							layout.Rigid(pg.inputRow(pg.addressEditor)),
							layout.Rigid(pg.inputRow(pg.signatureEditor)),
							layout.Rigid(pg.inputRow(pg.messageEditor)),
							layout.Rigid(pg.verifyAndClearButtons()),
							layout.Rigid(pg.verifyMessageResponse()),
						)
					})
				})
			},
			InfoTemplate: modal.VerifyMessageInfoTemplate,
		}
		return sp.Layout(gtx)
	}
	return components.UniformPadding(gtx, body)
}

func (pg *VerifyMessagePage) inputRow(editor decredmaterial.Editor) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, editor.Layout)
	}
}

func (pg *VerifyMessagePage) description() layout.Widget {
	return func(gtx C) D {
		desc := pg.Theme.Caption("Enter the address, signature, and message to verify:")
		desc.Color = pg.Theme.Color.GrayText2
		return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, desc.Layout)
	}
}

func (pg *VerifyMessagePage) verifyAndClearButtons() layout.Widget {
	return func(gtx C) D {
		dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, pg.clearBtn.Layout)
						}),
						layout.Rigid(pg.verifyButton.Layout),
					)
				})
			}),
		)
		return dims
	}
}

func (pg *VerifyMessagePage) verifyMessageResponse() layout.Widget {
	return func(gtx C) D {
		if pg.verifyMessageStatus != nil {
			return layout.Inset{Top: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
				pg.Theme.Separator().Layout(gtx)
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return pg.verifyMessageStatus.Layout(gtx, values.MarginPadding20)
							})
						}),
						layout.Rigid(pg.verifyMessage.Layout),
					)
				})
			})
		}
		return D{}
	}
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *VerifyMessagePage) HandleUserInteractions() {
	pg.verifyButton.SetEnabled(pg.updateBtn())

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(pg.addressEditor.Editor, pg.messageEditor.Editor, pg.signatureEditor.Editor)
	if isChanged {
		if pg.addressEditor.Editor.Focused() {
			pg.validateAddress()
		}
		pg.clearMessages()
	}

	if (pg.verifyButton.Clicked() || isSubmit) && pg.validateAllInputs() {
		pg.verifyMessage.Text = ""
		pg.verifyMessageStatus = nil
		valid, err := pg.WL.MultiWallet.VerifyMessage(pg.addressEditor.Editor.Text(), pg.messageEditor.Editor.Text(), pg.signatureEditor.Editor.Text())
		if err != nil {
			pg.verifyMessage.Text = "Error verifying message: " + err.Error()
			pg.verifyMessage.Color = pg.Theme.Color.Danger
			pg.verifyMessageStatus = decredmaterial.NewIcon(pg.Theme.Icons.NavigationCancel)
			pg.verifyMessageStatus.Color = pg.Theme.Color.Danger
			return
		}
		if !valid {
			pg.verifyMessage.Text = "Invalid signature or message"
			pg.verifyMessage.Color = pg.Theme.Color.Danger
			pg.verifyMessageStatus = decredmaterial.NewIcon(pg.Theme.Icons.NavigationCancel)
			pg.verifyMessageStatus.Color = pg.Theme.Color.Danger

			return
		}

		pg.verifyMessageStatus = decredmaterial.NewIcon(pg.Theme.Icons.ActionCheck)
		pg.verifyMessageStatus.Color = pg.Theme.Color.Success
		pg.verifyMessage.Text = "Valid signature"
		pg.verifyMessage.Color = pg.Theme.Color.Success
	}

	if pg.clearBtn.Clicked() {
		pg.clearInputs()
	}
}

// HandleKeyEvent is called when a key is pressed on the current window.
// Satisfies the load.KeyEventHandler interface for receiving key events.
func (pg *VerifyMessagePage) HandleKeyEvent(evt *key.Event) {
	// Switch editors on tab press.
	decredmaterial.SwitchEditors(evt, pg.addressEditor.Editor, pg.signatureEditor.Editor, pg.messageEditor.Editor)
}

func (pg *VerifyMessagePage) validateAllInputs() bool {
	if !pg.validateAddress() {
		return false
	}

	if !components.StringNotEmpty(pg.signatureEditor.Editor.Text()) {
		pg.signatureEditor.SetError("Field cannot be empty. Please provide valid signature.")
		return false
	}

	if !components.StringNotEmpty(pg.messageEditor.Editor.Text()) {
		pg.messageEditor.SetError("Field cannot be empty. Please provide valid signed message.")
		return false
	}

	return true
}

func (pg *VerifyMessagePage) updateBtn() bool {
	if pg.addressIsValid || components.StringNotEmpty(pg.signatureEditor.Editor.Text()) || components.StringNotEmpty(pg.messageEditor.Editor.Text()) {
		return true
	}
	return false
}

func (pg *VerifyMessagePage) clearInputs() {
	pg.verifyMessageStatus = nil
	pg.addressEditor.Editor.SetText("")
	pg.signatureEditor.Editor.SetText("")
	pg.messageEditor.Editor.SetText("")
	pg.verifyMessage.Text = ""
	pg.addressEditor.SetError("")
}

func (pg *VerifyMessagePage) clearMessages() {
	pg.verifyMessageStatus = nil
	pg.verifyMessage.Text = ""
}

func (pg *VerifyMessagePage) validateAddress() bool {
	address := pg.addressEditor.Editor.Text()
	pg.addressEditor.SetError("")

	var valid bool
	var errorMessage string

	switch {
	case !components.StringNotEmpty(address):
		errorMessage = "Please enter a valid address"
	case !pg.WL.MultiWallet.IsAddressValid(address):
		errorMessage = "Invalid address"
	default:
		valid = true
	}
	if !valid {
		pg.addressEditor.SetError(errorMessage)
	}

	pg.addressIsValid = valid
	return valid
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *VerifyMessagePage) OnNavigatedFrom() {}
