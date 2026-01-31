#include <windows.h>
#include <sphelper.h>
#include "tts.h"

ISpVoice *pVoice = NULL;

void initTTS() {
    CoInitialize(NULL);
    CoCreateInstance(CLSID_SpVoice, NULL, CLSCTX_ALL, IID_ISpVoice, (void **)&pVoice);

    IEnumSpObjectTokens* pEnum = NULL;
    ISpObjectToken* pToken = NULL;
    if (SUCCEEDED(SpEnumTokens(SPCAT_VOICES, NULL, NULL, &pEnum))) {
        bool foundChinese = false;

        while (pEnum->Next(1, &pToken, NULL) == S_OK) {
            LPWSTR pID = NULL;
            SpGetDescription(pToken, &pID);

            if (!foundChinese && (wcsstr(pID, L"Chinese") || wcsstr(pID, L"CHS"))) {
                pVoice->SetVoice(pToken);
                foundChinese = true;
            }

            CoTaskMemFree(pID);
            pToken->Release();
        }
        pEnum->Release();
    }
}

void speakText(const wchar_t *text) {
    pVoice->Speak(text, SVSFDefault, NULL);
}

void releaseTTS() {
    pVoice->Release();
    pVoice = NULL;
    CoUninitialize();
}
