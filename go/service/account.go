// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package service

import (
	"github.com/keybase/client/go/engine"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	"golang.org/x/net/context"
)

type AccountHandler struct {
	*BaseHandler
	libkb.Contextified
}

func NewAccountHandler(xp rpc.Transporter, g *libkb.GlobalContext) *AccountHandler {
	return &AccountHandler{
		BaseHandler:  NewBaseHandler(xp),
		Contextified: libkb.NewContextified(g),
	}
}

func (h *AccountHandler) PassphraseChange(_ context.Context, arg keybase1.PassphraseChangeArg) error {
	eng := engine.NewPassphraseChange(&arg, h.G())
	ctx := &engine.Context{
		SecretUI: h.getSecretUI(arg.SessionID),
	}
	return engine.RunEngine(eng, ctx)
}

func (h *AccountHandler) PassphrasePrompt(_ context.Context, arg keybase1.PassphrasePromptArg) (keybase1.GetPassphraseRes, error) {
	ui, err := h.G().UIRouter.GetSecretUI()
	if err != nil {
		return keybase1.GetPassphraseRes{}, err
	}
	if ui == nil {
		h.G().Log.Debug("delegate secret ui unavailable, falling back to standard one")
		ui = h.getSecretUI(arg.SessionID)
	}

	return ui.GetPassphrase(arg.GuiArg, nil)
}
