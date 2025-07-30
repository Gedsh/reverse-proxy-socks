package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getHandler(transport *http.Transport) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			addOptionsCORSHeaders(w)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		targetRaw := r.URL.Path[1:]
		targetURLStr, err := url.QueryUnescape(targetRaw)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusBadRequest)
			return
		}

		targetURL, err := url.Parse(targetURLStr)
		if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
			http.Error(w, "Malformed target URL", http.StatusBadRequest)
			return
		}

		outReq := new(http.Request)
		*outReq = *r
		outReq.URL = targetURL
		outReq.RequestURI = ""

		resp, err := transport.RoundTrip(outReq)
		if err != nil {
			http.Error(w, "Proxy error: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, vv := range resp.Header {
			if strings.ToLower(k) == "content-length" ||
				strings.HasPrefix(strings.ToLower(k), "access-control-") {
				continue
			}
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)

		if strings.HasSuffix(strings.ToLower(targetURL.Path), ".m3u8") {
			rewriteM3U8(w, resp.Body, targetURL)
		} else {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				log.Println("Copy error:", err)
			}
		}
	})
}

func rewriteM3U8(w http.ResponseWriter, body io.Reader, baseURL *url.URL) {
	scanner := bufio.NewScanner(body)
	var buf bytes.Buffer
	var pendingStreamTag string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
			pendingStreamTag = line
		} else if pendingStreamTag != "" {
			ref, err := baseURL.Parse(line)
			if err == nil {
				abs := ref.String()
				buf.WriteString(pendingStreamTag + "\n")
				buf.WriteString(abs + "\n")
			} else {
				buf.WriteString(pendingStreamTag + "\n")
				buf.WriteString(line + "\n")
			}
			pendingStreamTag = ""
		} else if strings.HasPrefix(line, "#EXT-X-MAP:") {
			const uriPrefix = `URI="`
			if idx := strings.Index(line, uriPrefix); idx != -1 {
				start := idx + len(uriPrefix)
				end := strings.Index(line[start:], `"`)
				if end != -1 {
					uri := line[start : start+end]
					ref, err := baseURL.Parse(uri)
					if err == nil {
						abs := ref.String()
						line = strings.Replace(line, uri, abs, 1)
					}
				}
			}
			buf.WriteString(line + "\n")
		} else if strings.HasPrefix(line, "#") {
			buf.WriteString(line + "\n")
		} else {
			ref, err := baseURL.Parse(line)
			if err == nil && ref.IsAbs() {
				buf.WriteString(ref.String() + "\n")
			} else {
				buf.WriteString(line + "\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("scanner error:", err)
	}
	_, err := w.Write(buf.Bytes())
	if err != nil {
		log.Println("write error:", err)
	}

}

func addOptionsCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, PATCH, HEAD, CONNECT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Range, X-Requested-With")
}
