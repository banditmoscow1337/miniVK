//go:build darwin

package binding

import "testing"

func TestDarwinKeyAndModifierTranslation(t *testing.T) {
	cases := []struct {
		macCode uint16
		want    int
	}{
		{0x00, GLFW_KEY_A},
		{0x01, GLFW_KEY_S},
		{0x02, GLFW_KEY_D},
		{0x03, GLFW_KEY_F},
		{0x04, GLFW_KEY_H},
		{0x06, GLFW_KEY_Z},
		{0x08, GLFW_KEY_C},
		{0x0C, GLFW_KEY_Q},
		{0x0D, GLFW_KEY_W},
		{0x0E, GLFW_KEY_E},
		{0x0F, GLFW_KEY_R},
		{0x10, GLFW_KEY_Y},
		{0x12, GLFW_KEY_1},
		{0x16, GLFW_KEY_6},
		{0x17, GLFW_KEY_5},
		{0x1D, GLFW_KEY_0},
		{0x1E, GLFW_KEY_RIGHT_BRACKET},
		{0x23, GLFW_KEY_P},
		{0x24, GLFW_KEY_ENTER},
		{0x27, GLFW_KEY_APOSTROPHE},
		{0x2D, GLFW_KEY_N},
		{0x30, GLFW_KEY_TAB},
		{0x31, GLFW_KEY_SPACE},
		{0x32, GLFW_KEY_GRAVE_ACCENT},
		{0x33, GLFW_KEY_BACKSPACE},
		{0x35, GLFW_KEY_ESCAPE},
		{0x39, GLFW_KEY_CAPS_LOCK},
		{0x3B, GLFW_KEY_LEFT_CONTROL},
		{0x3F, GLFW_KEY_UNKNOWN},
		{0x47, GLFW_KEY_NUM_LOCK},
		{0x60, GLFW_KEY_F5},
		{0x61, GLFW_KEY_F6},
		{0x64, GLFW_KEY_F8},
		{0x69, GLFW_KEY_PRINT_SCREEN},
		{0x6E, GLFW_KEY_MENU},
		{0x6F, GLFW_KEY_F12},
		{0x72, GLFW_KEY_INSERT},
		{0x73, GLFW_KEY_HOME},
		{0x74, GLFW_KEY_PAGE_UP},
		{0x75, GLFW_KEY_DELETE},
		{0x77, GLFW_KEY_END},
		{0x79, GLFW_KEY_PAGE_DOWN},
		{0x7A, GLFW_KEY_F1},
		{0x7B, GLFW_KEY_LEFT},
		{0x7C, GLFW_KEY_RIGHT},
		{0x7D, GLFW_KEY_DOWN},
		{0x7E, GLFW_KEY_UP},
		{0xFF, GLFW_KEY_UNKNOWN},
		{0x100, GLFW_KEY_UNKNOWN},
	}
	for _, tc := range cases {
		if got := translateMacKey(tc.macCode); got != tc.want {
			t.Errorf("translateMacKey(0x%X) = %d, want %d", tc.macCode, got, tc.want)
		}
	}

	flags := uintptr(nsEventModifierFlagShift | nsEventModifierFlagControl | nsEventModifierFlagOption | nsEventModifierFlagCommand | nsEventModifierFlagCapsLock)
	mods := cocoaMods(flags)
	wantMods := GLFW_MOD_SHIFT | GLFW_MOD_CONTROL | GLFW_MOD_ALT | GLFW_MOD_SUPER | GLFW_MOD_CAPS_LOCK
	if mods != wantMods {
		t.Fatalf("cocoaMods(0x%X) = 0x%X, want 0x%X", flags, mods, wantMods)
	}
}
