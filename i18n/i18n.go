package i18n

import (
	acceptlang "github.com/kazegusuri/grpc-accept-language"
	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var defaultLanguage = "en"

func SetDefaultLanguage(lang string) {
	defaultLanguage = lang
}

var _ grpc.UnaryServerInterceptor = UnaryI18nHandler

type tfuncKey struct{}

func UnaryI18nHandler(origctx context.Context, origreq interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return acceptlang.UnaryAcceptLanguageHandler(origctx, origreq, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		acceptLangs := acceptlang.FromContext(ctx)
		langs := acceptLangs.Languages()
		langs = append(langs, defaultLanguage)
		tfunc := i18n.MustTfunc(langs[0], langs[1:]...)
		ctx = context.WithValue(ctx, tfuncKey{}, tfunc)
		return handler(ctx, req)
	})
}

func MustTfunc(ctx context.Context) i18n.TranslateFunc {
	tfunc, ok := ctx.Value(tfuncKey{}).(i18n.TranslateFunc)
	if !ok {
		panic("could not find TranslateFunc from context")
	}
	return tfunc
}
