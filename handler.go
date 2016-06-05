package acceptlang

import (
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _ grpc.UnaryServerInterceptor = UnaryAcceptLanguageHandler

type alKey struct{}

func UnaryAcceptLanguageHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	acceptLangs := HandleAcceptLanguage(ctx)
	ctx = context.WithValue(ctx, alKey{}, acceptLangs)
	return handler(ctx, req)
}

func FromContext(ctx context.Context) AcceptLanguages {
	al, ok := ctx.Value(alKey{}).(AcceptLanguages)
	if !ok || al == nil {
		return []AcceptLanguage{}
	}
	return al
}

func HandleAcceptLanguage(ctx context.Context) AcceptLanguages {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil
	}

	header, ok := md["accept-language"]
	if !ok || len(header) == 0 {
		return nil
	}

	acceptLangHeader := header[0]
	acceptLangHeaderSlice := strings.Split(acceptLangHeader, ",")

	acceptLangs := make(AcceptLanguages, len(acceptLangHeaderSlice))
	for i, lang := range acceptLangHeaderSlice {
		lang = strings.TrimSpace(lang)
		qualSlice := strings.Split(lang, ";q=")
		if len(qualSlice) == 2 {
			qual, err := strconv.ParseFloat(qualSlice[1], 32)
			if err != nil {
				acceptLangs[i] = newAcceptLanguage(qualSlice[0], 1)
			} else {
				acceptLangs[i] = newAcceptLanguage(qualSlice[0], float32(qual))
			}
		} else {
			acceptLangs[i] = newAcceptLanguage(lang, 1)
		}
	}

	sort.Sort(sort.Reverse(acceptLangs))
	return acceptLangs
}

type AcceptLanguage struct {
	Language string
	Quality  float32
}

func newAcceptLanguage(lang string, qual float32) AcceptLanguage {
	return AcceptLanguage{Language: lang, Quality: qual}
}

type AcceptLanguages []AcceptLanguage

func (al AcceptLanguages) Languages() []string {
	langs := make([]string, len(al))
	for i := range al {
		langs[i] = al[i].Language
	}
	return langs
}

func (al AcceptLanguages) Len() int {
	return len(al)
}

func (al AcceptLanguages) Swap(i, j int) {
	al[i], al[j] = al[j], al[i]
}

func (al AcceptLanguages) Less(i, j int) bool {
	return al[i].Quality < al[j].Quality
}
