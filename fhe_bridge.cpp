#include "fhe_bridge.h"
#include "seal/seal.h"
#include <string>
#include <cstring>
#include <sstream>
#include <vector>
#include <iomanip>
#include <iostream>

using namespace seal;
using namespace std;

struct FHEContext {
    shared_ptr<SEALContext> context;
    shared_ptr<KeyGenerator> keygen;
    shared_ptr<PublicKey> public_key;
    shared_ptr<SecretKey> secret_key;
    shared_ptr<RelinKeys> relin_keys;
    shared_ptr<Encryptor> encryptor;
    shared_ptr<Evaluator> evaluator;
    shared_ptr<Decryptor> decryptor;
    shared_ptr<BatchEncoder> encoder;
    bool initialized;
    
    FHEContext() : initialized(false) {}
};

string serialize_ct(shared_ptr<SEALContext> context, Ciphertext& ct) {
    stringstream ss;
    ct.save(ss);
    string data = ss.str();
    stringstream hex_ss;
    hex_ss << hex << uppercase;
    for (unsigned char c : data) {
        hex_ss << setw(2) << setfill('0') << (int)c;
    }
    return hex_ss.str();
}

Ciphertext deserialize_ct(shared_ptr<SEALContext> context, const string& hex_str) {
    string data;
    for (size_t i = 0; i < hex_str.length(); i += 2) {
        string byte_str = hex_str.substr(i, 2);
        char byte = (char)strtol(byte_str.c_str(), nullptr, 16);
        data += byte;
    }
    stringstream ss(data);
    Ciphertext ct;
    ct.load(*context, ss);
    return ct;
}

extern "C" {

void* fhe_init() {
    FHEContext* ctx = new FHEContext();
    
    try {
        EncryptionParameters parms(scheme_type::bfv);
        parms.set_poly_modulus_degree(8192);
        parms.set_coeff_modulus(CoeffModulus::BFVDefault(8192));
        parms.set_plain_modulus(PlainModulus::Batching(8192, 20));
        
        ctx->context = make_shared<SEALContext>(parms);
        ctx->keygen = make_shared<KeyGenerator>(*ctx->context);
        ctx->secret_key = make_shared<SecretKey>(ctx->keygen->secret_key());
        ctx->public_key = make_shared<PublicKey>();
        ctx->keygen->create_public_key(*ctx->public_key);
        ctx->relin_keys = make_shared<RelinKeys>();
        ctx->keygen->create_relin_keys(*ctx->relin_keys);
        
        ctx->encryptor = make_shared<Encryptor>(*ctx->context, *ctx->public_key);
        ctx->evaluator = make_shared<Evaluator>(*ctx->context);
        ctx->decryptor = make_shared<Decryptor>(*ctx->context, *ctx->secret_key);
        ctx->encoder = make_shared<BatchEncoder>(*ctx->context);
        
        ctx->initialized = true;
    } catch (const exception& e) {
        cerr << "FHE_INIT ERROR: " << e.what() << endl;
    }
    
    return ctx;
}

char* fhe_encrypt(void* ctx_ptr, long value) {
    FHEContext* ctx = (FHEContext*)ctx_ptr;
    if (!ctx || !ctx->initialized) {
        cerr << "FHE_ENCRYPT: context not initialized" << endl;
        return nullptr;
    }
    
    try {
        vector<uint64_t> values(ctx->encoder->slot_count(), (uint64_t)value);
        Plaintext pt;
        ctx->encoder->encode(values, pt);
        
        Ciphertext ct(*ctx->context);
        ctx->encryptor->encrypt(pt, ct);
        
        string hex_str = serialize_ct(ctx->context, ct);
        char* result = (char*)malloc(hex_str.length() + 1);
        strcpy(result, hex_str.c_str());
        return result;
    } catch (const exception& e) {
        cerr << "FHE_ENCRYPT ERROR: " << e.what() << endl;
        return nullptr;
    }
}

long fhe_decrypt(void* ctx_ptr, const char* encrypted) {
    FHEContext* ctx = (FHEContext*)ctx_ptr;
    if (!ctx || !ctx->initialized || !encrypted) return 0;
    
    try {
        Ciphertext ct = deserialize_ct(ctx->context, string(encrypted));
        Plaintext pt;
        ctx->decryptor->decrypt(ct, pt);
        vector<uint64_t> values;
        ctx->encoder->decode(pt, values);
        return values.empty() ? 0 : (long)values[0];
    } catch (const exception& e) {
        cerr << "FHE_DECRYPT ERROR: " << e.what() << endl;
        return 0;
    }
}

char* fhe_add(void* ctx_ptr, const char* enc1, const char* enc2) {
    FHEContext* ctx = (FHEContext*)ctx_ptr;
    if (!ctx || !ctx->initialized || !enc1 || !enc2) return nullptr;
    
    try {
        Ciphertext ct1 = deserialize_ct(ctx->context, string(enc1));
        Ciphertext ct2 = deserialize_ct(ctx->context, string(enc2));
        ctx->evaluator->add_inplace(ct1, ct2);
        
        string hex_str = serialize_ct(ctx->context, ct1);
        char* result = (char*)malloc(hex_str.length() + 1);
        strcpy(result, hex_str.c_str());
        return result;
    } catch (...) { return nullptr; }
}

char* fhe_multiply(void* ctx_ptr, const char* enc1, const char* enc2) {
    FHEContext* ctx = (FHEContext*)ctx_ptr;
    if (!ctx || !ctx->initialized || !enc1 || !enc2) return nullptr;
    
    try {
        Ciphertext ct1 = deserialize_ct(ctx->context, string(enc1));
        Ciphertext ct2 = deserialize_ct(ctx->context, string(enc2));
        ctx->evaluator->multiply_inplace(ct1, ct2);
        ctx->evaluator->relinearize_inplace(ct1, *ctx->relin_keys);
        
        string hex_str = serialize_ct(ctx->context, ct1);
        char* result = (char*)malloc(hex_str.length() + 1);
        strcpy(result, hex_str.c_str());
        return result;
    } catch (...) { return nullptr; }
}

int fhe_noise_budget(void* ctx_ptr, const char* encrypted) {
    FHEContext* ctx = (FHEContext*)ctx_ptr;
    if (!ctx || !ctx->initialized || !encrypted) return -1;
    
    try {
        Ciphertext ct = deserialize_ct(ctx->context, string(encrypted));
        return ctx->decryptor->invariant_noise_budget(ct);
    } catch (...) { return -1; }
}

char* fhe_zans_add(void* ctx_ptr, const char* encrypted, int iterations) {
    FHEContext* ctx = (FHEContext*)ctx_ptr;
    if (!ctx || !ctx->initialized || !encrypted) return nullptr;
    
    try {
        Ciphertext ct = deserialize_ct(ctx->context, string(encrypted));
        vector<uint64_t> zero_values(ctx->encoder->slot_count(), 0ULL);
        Plaintext pt_zero;
        ctx->encoder->encode(zero_values, pt_zero);
        Ciphertext ct_zero(*ctx->context);
        
        for (int i = 0; i < iterations; i++) {
            ctx->encryptor->encrypt(pt_zero, ct_zero);
            ctx->evaluator->add_inplace(ct, ct_zero);
        }
        
        string hex_str = serialize_ct(ctx->context, ct);
        char* result = (char*)malloc(hex_str.length() + 1);
        strcpy(result, hex_str.c_str());
        return result;
    } catch (...) { return nullptr; }
}

void fhe_free(void* ctx_ptr) {
    delete (FHEContext*)ctx_ptr;
}

void fhe_free_string(char* str) {
    free(str);
}

}
