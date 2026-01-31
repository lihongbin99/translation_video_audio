#ifndef TTS_H
#define TTS_H
#define EXPORTS_API extern "C" __declspec(dllexport)

EXPORTS_API void initTTS();

EXPORTS_API void speakText(const wchar_t *text);

EXPORTS_API void releaseTTS();

#endif