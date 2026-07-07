#ifndef FHE_C_H
#define FHE_C_H

#ifdef __cplusplus
extern "C" {
#endif

void* fhe_init();
void fhe_generate_keys(void* ctx);
char* fhe_encrypt(void* ctx, long value);
long fhe_decrypt(void* ctx, const char* encrypted);
char* fhe_add(void* ctx, const char* enc1, const char* enc2);
char* fhe_multiply(void* ctx, const char* enc1, const char* enc2);
int fhe_noise_budget(void* ctx, const char* encrypted);
char* fhe_zans_add(void* ctx, const char* encrypted, int iterations);
void fhe_free(void* ctx);
void fhe_free_string(char* str);

#ifdef __cplusplus
}
#endif

#endif
