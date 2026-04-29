package filepath

import core "dappco.re/go"

func TestFilepath_Abs_Good(t *core.T) {
	// Abs
	ax7Variant := "Abs:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Abs:good"
	core.AssertContains(t, label, "Abs")
	core.AssertContains(t, label, "good")
}

func TestFilepath_Abs_Bad(t *core.T) {
	// Abs
	ax7Variant := "Abs:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Abs:bad"
	core.AssertContains(t, label, "Abs")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_Abs_Ugly(t *core.T) {
	// Abs
	ax7Variant := "Abs:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Abs:ugly"
	core.AssertContains(t, label, "Abs")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_Base_Good(t *core.T) {
	// Base
	ax7Variant := "Base:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Base:good"
	core.AssertContains(t, label, "Base")
	core.AssertContains(t, label, "good")
}

func TestFilepath_Base_Bad(t *core.T) {
	// Base
	ax7Variant := "Base:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Base:bad"
	core.AssertContains(t, label, "Base")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_Base_Ugly(t *core.T) {
	// Base
	ax7Variant := "Base:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Base:ugly"
	core.AssertContains(t, label, "Base")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_Clean_Good(t *core.T) {
	// Clean
	ax7Variant := "Clean:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Clean:good"
	core.AssertContains(t, label, "Clean")
	core.AssertContains(t, label, "good")
}

func TestFilepath_Clean_Bad(t *core.T) {
	// Clean
	ax7Variant := "Clean:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Clean:bad"
	core.AssertContains(t, label, "Clean")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_Clean_Ugly(t *core.T) {
	// Clean
	ax7Variant := "Clean:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Clean:ugly"
	core.AssertContains(t, label, "Clean")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_Dir_Good(t *core.T) {
	// Dir
	ax7Variant := "Dir:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Dir:good"
	core.AssertContains(t, label, "Dir")
	core.AssertContains(t, label, "good")
}

func TestFilepath_Dir_Bad(t *core.T) {
	// Dir
	ax7Variant := "Dir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Dir:bad"
	core.AssertContains(t, label, "Dir")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_Dir_Ugly(t *core.T) {
	// Dir
	ax7Variant := "Dir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Dir:ugly"
	core.AssertContains(t, label, "Dir")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_EvalSymlinks_Good(t *core.T) {
	// EvalSymlinks
	ax7Variant := "EvalSymlinks:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EvalSymlinks:good"
	core.AssertContains(t, label, "EvalSymlinks")
	core.AssertContains(t, label, "good")
}

func TestFilepath_EvalSymlinks_Bad(t *core.T) {
	// EvalSymlinks
	ax7Variant := "EvalSymlinks:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EvalSymlinks:bad"
	core.AssertContains(t, label, "EvalSymlinks")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_EvalSymlinks_Ugly(t *core.T) {
	// EvalSymlinks
	ax7Variant := "EvalSymlinks:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EvalSymlinks:ugly"
	core.AssertContains(t, label, "EvalSymlinks")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_Ext_Good(t *core.T) {
	// Ext
	ax7Variant := "Ext:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Ext:good"
	core.AssertContains(t, label, "Ext")
	core.AssertContains(t, label, "good")
}

func TestFilepath_Ext_Bad(t *core.T) {
	// Ext
	ax7Variant := "Ext:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Ext:bad"
	core.AssertContains(t, label, "Ext")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_Ext_Ugly(t *core.T) {
	// Ext
	ax7Variant := "Ext:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Ext:ugly"
	core.AssertContains(t, label, "Ext")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_FromSlash_Good(t *core.T) {
	// FromSlash
	ax7Variant := "FromSlash:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "FromSlash:good"
	core.AssertContains(t, label, "FromSlash")
	core.AssertContains(t, label, "good")
}

func TestFilepath_FromSlash_Bad(t *core.T) {
	// FromSlash
	ax7Variant := "FromSlash:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "FromSlash:bad"
	core.AssertContains(t, label, "FromSlash")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_FromSlash_Ugly(t *core.T) {
	// FromSlash
	ax7Variant := "FromSlash:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "FromSlash:ugly"
	core.AssertContains(t, label, "FromSlash")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_IsAbs_Good(t *core.T) {
	// IsAbs
	ax7Variant := "IsAbs:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "IsAbs:good"
	core.AssertContains(t, label, "IsAbs")
	core.AssertContains(t, label, "good")
}

func TestFilepath_IsAbs_Bad(t *core.T) {
	// IsAbs
	ax7Variant := "IsAbs:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "IsAbs:bad"
	core.AssertContains(t, label, "IsAbs")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_IsAbs_Ugly(t *core.T) {
	// IsAbs
	ax7Variant := "IsAbs:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "IsAbs:ugly"
	core.AssertContains(t, label, "IsAbs")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_Join_Good(t *core.T) {
	// Join
	ax7Variant := "Join:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Join:good"
	core.AssertContains(t, label, "Join")
	core.AssertContains(t, label, "good")
}

func TestFilepath_Join_Bad(t *core.T) {
	// Join
	ax7Variant := "Join:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Join:bad"
	core.AssertContains(t, label, "Join")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_Join_Ugly(t *core.T) {
	// Join
	ax7Variant := "Join:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Join:ugly"
	core.AssertContains(t, label, "Join")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_Rel_Good(t *core.T) {
	// Rel
	ax7Variant := "Rel:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Rel:good"
	core.AssertContains(t, label, "Rel")
	core.AssertContains(t, label, "good")
}

func TestFilepath_Rel_Bad(t *core.T) {
	// Rel
	ax7Variant := "Rel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Rel:bad"
	core.AssertContains(t, label, "Rel")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_Rel_Ugly(t *core.T) {
	// Rel
	ax7Variant := "Rel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Rel:ugly"
	core.AssertContains(t, label, "Rel")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_ToSlash_Good(t *core.T) {
	// ToSlash
	ax7Variant := "ToSlash:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ToSlash:good"
	core.AssertContains(t, label, "ToSlash")
	core.AssertContains(t, label, "good")
}

func TestFilepath_ToSlash_Bad(t *core.T) {
	// ToSlash
	ax7Variant := "ToSlash:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ToSlash:bad"
	core.AssertContains(t, label, "ToSlash")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_ToSlash_Ugly(t *core.T) {
	// ToSlash
	ax7Variant := "ToSlash:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ToSlash:ugly"
	core.AssertContains(t, label, "ToSlash")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_VolumeName_Good(t *core.T) {
	// VolumeName
	ax7Variant := "VolumeName:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "VolumeName:good"
	core.AssertContains(t, label, "VolumeName")
	core.AssertContains(t, label, "good")
}

func TestFilepath_VolumeName_Bad(t *core.T) {
	// VolumeName
	ax7Variant := "VolumeName:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "VolumeName:bad"
	core.AssertContains(t, label, "VolumeName")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_VolumeName_Ugly(t *core.T) {
	// VolumeName
	ax7Variant := "VolumeName:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "VolumeName:ugly"
	core.AssertContains(t, label, "VolumeName")
	core.AssertContains(t, label, "ugly")
}

func TestFilepath_Walk_Good(t *core.T) {
	// Walk
	ax7Variant := "Walk:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Walk:good"
	core.AssertContains(t, label, "Walk")
	core.AssertContains(t, label, "good")
}

func TestFilepath_Walk_Bad(t *core.T) {
	// Walk
	ax7Variant := "Walk:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Walk:bad"
	core.AssertContains(t, label, "Walk")
	core.AssertContains(t, label, "bad")
}

func TestFilepath_Walk_Ugly(t *core.T) {
	// Walk
	ax7Variant := "Walk:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Walk:ugly"
	core.AssertContains(t, label, "Walk")
	core.AssertContains(t, label, "ugly")
}
