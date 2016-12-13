package goselenium

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// WindowHandleResponse is the response returned from the WindowHandle() method.
// The handle is the current active window. Should you switch windows,
// any value returned prior to that call will be invalid.
type WindowHandleResponse struct {
	State  string
	Handle string
}

// CloseWindowResponse is the response returned from the CloseWindow() method.
// As per the W3C specification, it yields all of the available window handles
// minus the active one that closes as a result of the CloseWindow() call.
type CloseWindowResponse struct {
	State   string   `json:state"`
	Handles []string `json:"value"`
}

// SwitchToWindowResponse is the response returned from the SwitchToWindow()
// method. You can verify that this result is correct by calling the
// WindowHandle() method. The two should match.
type SwitchToWindowResponse struct {
}

// WindowHandlesResponse is the response returned from the WindowHandles()
// method. This is essentially an array of available window handlers that
// aren't neccessarily active.
type WindowHandlesResponse struct {
	State   string   `json:"state"`
	Handles []string `json:"value"`
}

// SwitchToFrameResponse is the response returned from the SwitchToFrame()
// method. For now, according to the specification, it only returns a state.
type SwitchToFrameResponse struct {
	State string
}

// SwitchToParentFrameResponse represents the response from attempting to
// switch the top level browsing context to the parent of the current top level
// browsing context.
type SwitchToParentFrameResponse struct {
	State string
}

// WindowSizeResponse is the response returned from calling the WindowSize
// method. The definitions are in CSS pixels.
type WindowSizeResponse struct {
	State      string     `json:"state"`
	Dimensions Dimensions `json:"value"`
}

// Dimensions is a type that is both returned and accept by functions. It is
// usually only used for the window size components.
type Dimensions struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

// SetWindowSizeResponse is the response that is returned from setting the
// window size of the current top level browsing context.
type SetWindowSizeResponse struct {
	State string
}

// MaximizeWindowResponse is the response that is returned from increasing the
// browser to match the viewport.
type MaximizeWindowResponse struct {
	State string
}

func (s *seleniumWebDriver) WindowHandle() (*WindowHandleResponse, error) {
	var response WindowHandleResponse
	var err error

	if s.sessionID == "" {
		return nil, newSessionIDError("WindowHandle()")
	}

	url := fmt.Sprintf("%s/session/%s/window", s.seleniumURL, s.sessionID)

	resp, err := s.valueRequest(&request{
		url:           url,
		method:        "GET",
		body:          nil,
		callingMethod: "WindowHandle",
	})
	if err != nil {
		return nil, err
	}

	response = WindowHandleResponse{
		State:  resp.State,
		Handle: resp.Value,
	}

	return &response, nil
}

func (s *seleniumWebDriver) CloseWindow() (*CloseWindowResponse, error) {
	var response CloseWindowResponse
	var err error

	if s.sessionID == "" {
		return nil, newSessionIDError("CloseWindow()")
	}

	url := fmt.Sprintf("%s/session/%s/window", s.seleniumURL, s.sessionID)

	resp, err := s.apiService.performRequest(url, "DELETE", nil)
	if err != nil {
		return nil, newCommunicationError(err, "CloseWindow()", url, resp)
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, newUnmarshallingError(err, "CloseWindow()", string(resp))
	}

	return &response, nil
}

func (s *seleniumWebDriver) SwitchToWindow(handle string) (*SwitchToWindowResponse, error) {
	return nil, nil
}

func (s *seleniumWebDriver) WindowHandles() (*WindowHandlesResponse, error) {
	var response WindowHandlesResponse
	var err error

	if s.sessionID == "" {
		return nil, newSessionIDError("WindowHandles()")
	}

	url := fmt.Sprintf("%s/session/%s/window/handles", s.seleniumURL, s.sessionID)

	resp, err := s.apiService.performRequest(url, "GET", nil)
	if err != nil {
		return nil, newCommunicationError(err, "WindowHandles()", url, resp)
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, newUnmarshallingError(err, "WindowHandles()", string(resp))
	}

	return &response, nil
}

func (s *seleniumWebDriver) SwitchToFrame(by By) (*SwitchToFrameResponse, error) {
	var err error

	if s.sessionID == "" {
		return nil, newSessionIDError("SwitchToFrame()")
	}
	if by == nil || (by.Type() != "index") {
		return nil, newInvalidArgumentError("By was not of required type in SwitchToFrame()", "by", "")
	}

	url := fmt.Sprintf("%s/session/%s/frame", s.seleniumURL, s.sessionID)

	params := map[string]interface{}{
		"id": by.Value(),
	}
	requestJSON, err := json.Marshal(params)
	if err != nil {
		return nil, newMarshallingError(err, "SwitchToFrame()", params)
	}

	body := bytes.NewReader(requestJSON)
	resp, err := s.stateRequest(&request{
		url:           url,
		method:        "POST",
		body:          body,
		callingMethod: "SwitchToFrame()",
	})
	if err != nil {
		return nil, err
	}

	return &SwitchToFrameResponse{State: resp.State}, nil
}

func (s *seleniumWebDriver) SwitchToParentFrame() (*SwitchToParentFrameResponse, error) {
	var err error

	if s.sessionID == "" {
		return nil, newSessionIDError("SwitchToParentFrame()")
	}

	url := fmt.Sprintf("%s/session/%s/frame/parent", s.seleniumURL, s.sessionID)

	resp, err := s.stateRequest(&request{
		url:           url,
		method:        "POST",
		body:          nil,
		callingMethod: "SwitchToParentFrame()",
	})
	if err != nil {
		return nil, err
	}

	return &SwitchToParentFrameResponse{State: resp.State}, nil
}

func (s *seleniumWebDriver) WindowSize() (*WindowSizeResponse, error) {
	var response WindowSizeResponse
	var err error

	if s.sessionID == "" {
		return nil, newSessionIDError("WindowSize()")
	}

	url := fmt.Sprintf("%s/session/%s/window/size", s.seleniumURL, s.sessionID)

	resp, err := s.apiService.performRequest(url, "GET", nil)
	if err != nil {
		return nil, newCommunicationError(err, "WindowSize()", url, nil)
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, newUnmarshallingError(err, "WindowSize()", string(resp))
	}

	return &response, nil
}

func (s *seleniumWebDriver) SetWindowSize(dimension *Dimensions) (*SetWindowSizeResponse, error) {
	var err error

	if dimension == nil {
		return nil, newInvalidArgumentError("Dimension was nil", "dimension", "")
	} else if s.sessionID == "" {
		return nil, newSessionIDError("SetWindowSize()")
	}

	url := fmt.Sprintf("%s/session/%s/window/size", s.seleniumURL, s.sessionID)

	body := map[string]uint{
		"width":  dimension.Width,
		"height": dimension.Height,
	}
	json, err := json.Marshal(body)
	if err != nil {
		return nil, newMarshallingError(err, "SetWindowSize()", body)
	}

	jsonBytes := bytes.NewReader(json)
	resp, err := s.stateRequest(&request{
		url:           url,
		method:        "POST",
		body:          jsonBytes,
		callingMethod: "SetWindowSize()",
	})
	if err != nil {
		return nil, err
	}

	return &SetWindowSizeResponse{State: resp.State}, nil
}

func (s *seleniumWebDriver) MaximizeWindow() (*MaximizeWindowResponse, error) {
	var err error

	if s.sessionID == "" {
		return nil, newSessionIDError("MaximizeWindow()")
	}

	url := fmt.Sprintf("%s/session/%s/window/maximize", s.seleniumURL, s.sessionID)

	resp, err := s.stateRequest(&request{
		url:           url,
		method:        "POST",
		body:          nil,
		callingMethod: "MaximizeWindow()",
	})
	if err != nil {
		return nil, err
	}

	return &MaximizeWindowResponse{State: resp.State}, nil
}
