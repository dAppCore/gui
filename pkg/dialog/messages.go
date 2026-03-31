package dialog

type TaskOpenFile struct{ Options OpenFileOptions }

type TaskSaveFile struct{ Options SaveFileOptions }

type TaskOpenDirectory struct{ Options OpenDirectoryOptions }

type TaskMessageDialog struct{ Options MessageDialogOptions }
