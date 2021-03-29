package sDocument

import (
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"gopkg.in/gographics/imagick.v3/imagick"
	"strings"
)

/*
	ConvertImageFormat writes out the image in the same/different format

	EXAMPLE:
		if dsm, soteErr := New(); soteErr.ErrCode == nil {
			if _, soteErr = dsm.ConvertImageFormat("logo.jpg","out.png"); soteErr.ErrCode != nil {
				sLogger.Info(fmt.Sprintf("ConvertImageFormat failed:Expected soteErr to be nil:%v ", soteErr.FmtErrMsg))
			}
		}
*/
func SConvertImageFormat(sourcePath string, targetPath string) (pdfFilePtr *imagick.ImageCommandResult,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var serr error
	imagick.Initialize()
	defer imagick.Terminate() // Memory leak cleanup

	if pdfFilePtr, serr = imagick.ConvertImageCommand([]string{"convert", sourcePath, targetPath}); serr != nil {
		if strings.Contains(serr.Error(), "No such file or directory") {
			soteErr = sError.GetSError(199999, sError.BuildParams([]string{"Invalid upload path or filename"}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}
