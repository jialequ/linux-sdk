package internal

import (
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
	"github.com/jialequ/linux-sdk/core/hash"
	"github.com/stretchr/testify/assert"
)

const (
	b64pb                = `CpgBCgtoZWxsby5wcm90bxIFaGVsbG8iHQoHUmVxdWVzdBISCgRwaW5nGAEgASgJUgRwaW5nIh4KCFJlc3BvbnNlEhIKBHBvbmcYASABKAlSBHBvbmcyMAoFSGVsbG8SJwoEUGluZxIOLmhlbGxvLlJlcXVlc3QaDy5oZWxsby5SZXNwb25zZUIJWgcuL2hlbGxvYgZwcm90bzM=`
	b64pbWithAnnotations = `Cs4EChVnb29nbGUvYXBpL2h0dHAucHJvdG8SCmdvb2dsZS5hcGkieQoESHR0cBIqCgVydWxlcxgBIAMoCzIULmdvb2dsZS5hcGkuSHR0cFJ1bGVSBXJ1bGVzEkUKH2Z1bGx5X2RlY29kZV9yZXNlcnZlZF9leHBhbnNpb24YAiABKAhSHGZ1bGx5RGVjb2RlUmVzZXJ2ZWRFeHBhbnNpb24i2gIKCEh0dHBSdWxlEhoKCHNlbGVjdG9yGAEgASgJUghzZWxlY3RvchISCgNnZXQYAiABKAlIAFIDZ2V0EhIKA3B1dBgDIAEoCUgAUgNwdXQSFAoEcG9zdBgEIAEoCUgAUgRwb3N0EhgKBmRlbGV0ZRgFIAEoCUgAUgZkZWxldGUSFgoFcGF0Y2gYBiABKAlIAFIFcGF0Y2gSNwoGY3VzdG9tGAggASgLMh0uZ29vZ2xlLmFwaS5DdXN0b21IdHRwUGF0dGVybkgAUgZjdXN0b20SEgoEYm9keRgHIAEoCVIEYm9keRIjCg1yZXNwb25zZV9ib2R5GAwgASgJUgxyZXNwb25zZUJvZHkSRQoTYWRkaXRpb25hbF9iaW5kaW5ncxgLIAMoCzIULmdvb2dsZS5hcGkuSHR0cFJ1bGVSEmFkZGl0aW9uYWxCaW5kaW5nc0IJCgdwYXR0ZXJuIjsKEUN1c3RvbUh0dHBQYXR0ZXJuEhIKBGtpbmQYASABKAlSBGtpbmQSEgoEcGF0aBgCIAEoCVIEcGF0aEIMWgpnb29nbGUvYXBpYgZwcm90bzMK8zsKIGdvb2dsZS9wcm90b2J1Zi9kZXNjcmlwdG9yLnByb3RvEg9nb29nbGUucHJvdG9idWYiTQoRRmlsZURlc2NyaXB0b3JTZXQSOAoEZmlsZRgBIAMoCzIkLmdvb2dsZS5wcm90b2J1Zi5GaWxlRGVzY3JpcHRvclByb3RvUgRmaWxlIuQEChNGaWxlRGVzY3JpcHRvclByb3RvEhIKBG5hbWUYASABKAlSBG5hbWUSGAoHcGFja2FnZRgCIAEoCVIHcGFja2FnZRIeCgpkZXBlbmRlbmN5GAMgAygJUgpkZXBlbmRlbmN5EisKEXB1YmxpY19kZXBlbmRlbmN5GAogAygFUhBwdWJsaWNEZXBlbmRlbmN5EicKD3dlYWtfZGVwZW5kZW5jeRgLIAMoBVIOd2Vha0RlcGVuZGVuY3kSQwoMbWVzc2FnZV90eXBlGAQgAygLMiAuZ29vZ2xlLnByb3RvYnVmLkRlc2NyaXB0b3JQcm90b1ILbWVzc2FnZVR5cGUSQQoJZW51bV90eXBlGAUgAygLMiQuZ29vZ2xlLnByb3RvYnVmLkVudW1EZXNjcmlwdG9yUHJvdG9SCGVudW1UeXBlEkEKB3NlcnZpY2UYBiADKAsyJy5nb29nbGUucHJvdG9idWYuU2VydmljZURlc2NyaXB0b3JQcm90b1IHc2VydmljZRJDCglleHRlbnNpb24YByADKAsyJS5nb29nbGUucHJvdG9idWYuRmllbGREZXNjcmlwdG9yUHJvdG9SCWV4dGVuc2lvbhI2CgdvcHRpb25zGAggASgLMhwuZ29vZ2xlLnByb3RvYnVmLkZpbGVPcHRpb25zUgdvcHRpb25zEkkKEHNvdXJjZV9jb2RlX2luZm8YCSABKAsyHy5nb29nbGUucHJvdG9idWYuU291cmNlQ29kZUluZm9SDnNvdXJjZUNvZGVJbmZvEhYKBnN5bnRheBgMIAEoCVIGc3ludGF4IrkGCg9EZXNjcmlwdG9yUHJvdG8SEgoEbmFtZRgBIAEoCVIEbmFtZRI7CgVmaWVsZBgCIAMoCzIlLmdvb2dsZS5wcm90b2J1Zi5GaWVsZERlc2NyaXB0b3JQcm90b1IFZmllbGQSQwoJZXh0ZW5zaW9uGAYgAygLMiUuZ29vZ2xlLnByb3RvYnVmLkZpZWxkRGVzY3JpcHRvclByb3RvUglleHRlbnNpb24SQQoLbmVzdGVkX3R5cGUYAyADKAsyIC5nb29nbGUucHJvdG9idWYuRGVzY3JpcHRvclByb3RvUgpuZXN0ZWRUeXBlEkEKCWVudW1fdHlwZRgEIAMoCzIkLmdvb2dsZS5wcm90b2J1Zi5FbnVtRGVzY3JpcHRvclByb3RvUghlbnVtVHlwZRJYCg9leHRlbnNpb25fcmFuZ2UYBSADKAsyLy5nb29nbGUucHJvdG9idWYuRGVzY3JpcHRvclByb3RvLkV4dGVuc2lvblJhbmdlUg5leHRlbnNpb25SYW5nZRJECgpvbmVvZl9kZWNsGAggAygLMiUuZ29vZ2xlLnByb3RvYnVmLk9uZW9mRGVzY3JpcHRvclByb3RvUglvbmVvZkRlY2wSOQoHb3B0aW9ucxgHIAEoCzIfLmdvb2dsZS5wcm90b2J1Zi5NZXNzYWdlT3B0aW9uc1IHb3B0aW9ucxJVCg5yZXNlcnZlZF9yYW5nZRgJIAMoCzIuLmdvb2dsZS5wcm90b2J1Zi5EZXNjcmlwdG9yUHJvdG8uUmVzZXJ2ZWRSYW5nZVINcmVzZXJ2ZWRSYW5nZRIjCg1yZXNlcnZlZF9uYW1lGAogAygJUgxyZXNlcnZlZE5hbWUaegoORXh0ZW5zaW9uUmFuZ2USFAoFc3RhcnQYASABKAVSBXN0YXJ0EhAKA2VuZBgCIAEoBVIDZW5kEkAKB29wdGlvbnMYAyABKAsyJi5nb29nbGUucHJvdG9idWYuRXh0ZW5zaW9uUmFuZ2VPcHRpb25zUgdvcHRpb25zGjcKDVJlc2VydmVkUmFuZ2USFAoFc3RhcnQYASABKAVSBXN0YXJ0EhAKA2VuZBgCIAEoBVIDZW5kInwKFUV4dGVuc2lvblJhbmdlT3B0aW9ucxJYChR1bmludGVycHJldGVkX29wdGlvbhjnByADKAsyJC5nb29nbGUucHJvdG9idWYuVW5pbnRlcnByZXRlZE9wdGlvblITdW5pbnRlcnByZXRlZE9wdGlvbioJCOgHEICAgIACIsEGChRGaWVsZERlc2NyaXB0b3JQcm90bxISCgRuYW1lGAEgASgJUgRuYW1lEhYKBm51bWJlchgDIAEoBVIGbnVtYmVyEkEKBWxhYmVsGAQgASgOMisuZ29vZ2xlLnByb3RvYnVmLkZpZWxkRGVzY3JpcHRvclByb3RvLkxhYmVsUgVsYWJlbBI+CgR0eXBlGAUgASgOMiouZ29vZ2xlLnByb3RvYnVmLkZpZWxkRGVzY3JpcHRvclByb3RvLlR5cGVSBHR5cGUSGwoJdHlwZV9uYW1lGAYgASgJUgh0eXBlTmFtZRIaCghleHRlbmRlZRgCIAEoCVIIZXh0ZW5kZWUSIwoNZGVmYXVsdF92YWx1ZRgHIAEoCVIMZGVmYXVsdFZhbHVlEh8KC29uZW9mX2luZGV4GAkgASgFUgpvbmVvZkluZGV4EhsKCWpzb25fbmFtZRgKIAEoCVIIanNvbk5hbWUSNwoHb3B0aW9ucxgIIAEoCzIdLmdvb2dsZS5wcm90b2J1Zi5GaWVsZE9wdGlvbnNSB29wdGlvbnMSJwoPcHJvdG8zX29wdGlvbmFsGBEgASgIUg5wcm90bzNPcHRpb25hbCK2AgoEVHlwZRIPCgtUWVBFX0RPVUJMRRABEg4KClRZUEVfRkxPQVQQAhIOCgpUWVBFX0lOVDY0EAMSDwoLVFlQRV9VSU5UNjQQBBIOCgpUWVBFX0lOVDMyEAUSEAoMVFlQRV9GSVhFRDY0EAYSEAoMVFlQRV9GSVhFRDMyEAcSDQoJVFlQRV9CT09MEAgSDwoLVFlQRV9TVFJJTkcQCRIOCgpUWVBFX0dST1VQEAoSEAoMVFlQRV9NRVNTQUdFEAsSDgoKVFlQRV9CWVRFUxAMEg8KC1RZUEVfVUlOVDMyEA0SDQoJVFlQRV9FTlVNEA4SEQoNVFlQRV9TRklYRUQzMhAPEhEKDVRZUEVfU0ZJWEVENjQQEBIPCgtUWVBFX1NJTlQzMhAREg8KC1RZUEVfU0lOVDY0EBIiQwoFTGFiZWwSEgoOTEFCRUxfT1BUSU9OQUwQARISCg5MQUJFTF9SRVFVSVJFRBACEhIKDkxBQkVMX1JFUEVBVEVEEAMiYwoUT25lb2ZEZXNjcmlwdG9yUHJvdG8SEgoEbmFtZRgBIAEoCVIEbmFtZRI3CgdvcHRpb25zGAIgASgLMh0uZ29vZ2xlLnByb3RvYnVmLk9uZW9mT3B0aW9uc1IHb3B0aW9ucyLjAgoTRW51bURlc2NyaXB0b3JQcm90bxISCgRuYW1lGAEgASgJUgRuYW1lEj8KBXZhbHVlGAIgAygLMikuZ29vZ2xlLnByb3RvYnVmLkVudW1WYWx1ZURlc2NyaXB0b3JQcm90b1IFdmFsdWUSNgoHb3B0aW9ucxgDIAEoCzIcLmdvb2dsZS5wcm90b2J1Zi5FbnVtT3B0aW9uc1IHb3B0aW9ucxJdCg5yZXNlcnZlZF9yYW5nZRgEIAMoCzI2Lmdvb2dsZS5wcm90b2J1Zi5FbnVtRGVzY3JpcHRvclByb3RvLkVudW1SZXNlcnZlZFJhbmdlUg1yZXNlcnZlZFJhbmdlEiMKDXJlc2VydmVkX25hbWUYBSADKAlSDHJlc2VydmVkTmFtZRo7ChFFbnVtUmVzZXJ2ZWRSYW5nZRIUCgVzdGFydBgBIAEoBVIFc3RhcnQSEAoDZW5kGAIgASgFUgNlbmQigwEKGEVudW1WYWx1ZURlc2NyaXB0b3JQcm90bxISCgRuYW1lGAEgASgJUgRuYW1lEhYKBm51bWJlchgCIAEoBVIGbnVtYmVyEjsKB29wdGlvbnMYAyABKAsyIS5nb29nbGUucHJvdG9idWYuRW51bVZhbHVlT3B0aW9uc1IHb3B0aW9ucyKnAQoWU2VydmljZURlc2NyaXB0b3JQcm90bxISCgRuYW1lGAEgASgJUgRuYW1lEj4KBm1ldGhvZBgCIAMoCzImLmdvb2dsZS5wcm90b2J1Zi5NZXRob2REZXNjcmlwdG9yUHJvdG9SBm1ldGhvZBI5CgdvcHRpb25zGAMgASgLMh8uZ29vZ2xlLnByb3RvYnVmLlNlcnZpY2VPcHRpb25zUgdvcHRpb25zIokCChVNZXRob2REZXNjcmlwdG9yUHJvdG8SEgoEbmFtZRgBIAEoCVIEbmFtZRIdCgppbnB1dF90eXBlGAIgASgJUglpbnB1dFR5cGUSHwoLb3V0cHV0X3R5cGUYAyABKAlSCm91dHB1dFR5cGUSOAoHb3B0aW9ucxgEIAEoCzIeLmdvb2dsZS5wcm90b2J1Zi5NZXRob2RPcHRpb25zUgdvcHRpb25zEjAKEGNsaWVudF9zdHJlYW1pbmcYBSABKAg6BWZhbHNlUg9jbGllbnRTdHJlYW1pbmcSMAoQc2VydmVyX3N0cmVhbWluZxgGIAEoCDoFZmFsc2VSD3NlcnZlclN0cmVhbWluZyKRCQoLRmlsZU9wdGlvbnMSIQoMamF2YV9wYWNrYWdlGAEgASgJUgtqYXZhUGFja2FnZRIwChRqYXZhX291dGVyX2NsYXNzbmFtZRgIIAEoCVISamF2YU91dGVyQ2xhc3NuYW1lEjUKE2phdmFfbXVsdGlwbGVfZmlsZXMYCiABKAg6BWZhbHNlUhFqYXZhTXVsdGlwbGVGaWxlcxJECh1qYXZhX2dlbmVyYXRlX2VxdWFsc19hbmRfaGFzaBgUIAEoCEICGAFSGWphdmFHZW5lcmF0ZUVxdWFsc0FuZEhhc2gSOgoWamF2YV9zdHJpbmdfY2hlY2tfdXRmOBgbIAEoCDoFZmFsc2VSE2phdmFTdHJpbmdDaGVja1V0ZjgSUwoMb3B0aW1pemVfZm9yGAkgASgOMikuZ29vZ2xlLnByb3RvYnVmLkZpbGVPcHRpb25zLk9wdGltaXplTW9kZToFU1BFRURSC29wdGltaXplRm9yEh0KCmdvX3BhY2thZ2UYCyABKAlSCWdvUGFja2FnZRI1ChNjY19nZW5lcmljX3NlcnZpY2VzGBAgASgIOgVmYWxzZVIRY2NHZW5lcmljU2VydmljZXMSOQoVamF2YV9nZW5lcmljX3NlcnZpY2VzGBEgASgIOgVmYWxzZVITamF2YUdlbmVyaWNTZXJ2aWNlcxI1ChNweV9nZW5lcmljX3NlcnZpY2VzGBIgASgIOgVmYWxzZVIRcHlHZW5lcmljU2VydmljZXMSNwoUcGhwX2dlbmVyaWNfc2VydmljZXMYKiABKAg6BWZhbHNlUhJwaHBHZW5lcmljU2VydmljZXMSJQoKZGVwcmVjYXRlZBgXIAEoCDoFZmFsc2VSCmRlcHJlY2F0ZWQSLgoQY2NfZW5hYmxlX2FyZW5hcxgfIAEoCDoEdHJ1ZVIOY2NFbmFibGVBcmVuYXMSKgoRb2JqY19jbGFzc19wcmVmaXgYJCABKAlSD29iamNDbGFzc1ByZWZpeBIpChBjc2hhcnBfbmFtZXNwYWNlGCUgASgJUg9jc2hhcnBOYW1lc3BhY2USIQoMc3dpZnRfcHJlZml4GCcgASgJUgtzd2lmdFByZWZpeBIoChBwaHBfY2xhc3NfcHJlZml4GCggASgJUg5waHBDbGFzc1ByZWZpeBIjCg1waHBfbmFtZXNwYWNlGCkgASgJUgxwaHBOYW1lc3BhY2USNAoWcGhwX21ldGFkYXRhX25hbWVzcGFjZRgsIAEoCVIUcGhwTWV0YWRhdGFOYW1lc3BhY2USIQoMcnVieV9wYWNrYWdlGC0gASgJUgtydWJ5UGFja2FnZRJYChR1bmludGVycHJldGVkX29wdGlvbhjnByADKAsyJC5nb29nbGUucHJvdG9idWYuVW5pbnRlcnByZXRlZE9wdGlvblITdW5pbnRlcnByZXRlZE9wdGlvbiI6CgxPcHRpbWl6ZU1vZGUSCQoFU1BFRUQQARINCglDT0RFX1NJWkUQAhIQCgxMSVRFX1JVTlRJTUUQAyoJCOgHEICAgIACSgQIJhAnIuMCCg5NZXNzYWdlT3B0aW9ucxI8ChdtZXNzYWdlX3NldF93aXJlX2Zvcm1hdBgBIAEoCDoFZmFsc2VSFG1lc3NhZ2VTZXRXaXJlRm9ybWF0EkwKH25vX3N0YW5kYXJkX2Rlc2NyaXB0b3JfYWNjZXNzb3IYAiABKAg6BWZhbHNlUhxub1N0YW5kYXJkRGVzY3JpcHRvckFjY2Vzc29yEiUKCmRlcHJlY2F0ZWQYAyABKAg6BWZhbHNlUgpkZXByZWNhdGVkEhsKCW1hcF9lbnRyeRgHIAEoCFIIbWFwRW50cnkSWAoUdW5pbnRlcnByZXRlZF9vcHRpb24Y5wcgAygLMiQuZ29vZ2xlLnByb3RvYnVmLlVuaW50ZXJwcmV0ZWRPcHRpb25SE3VuaW50ZXJwcmV0ZWRPcHRpb24qCQjoBxCAgICAAkoECAQQBUoECAUQBkoECAYQB0oECAgQCUoECAkQCiKSBAoMRmllbGRPcHRpb25zEkEKBWN0eXBlGAEgASgOMiMuZ29vZ2xlLnByb3RvYnVmLkZpZWxkT3B0aW9ucy5DVHlwZToGU1RSSU5HUgVjdHlwZRIWCgZwYWNrZWQYAiABKAhSBnBhY2tlZBJHCgZqc3R5cGUYBiABKA4yJC5nb29nbGUucHJvdG9idWYuRmllbGRPcHRpb25zLkpTVHlwZToJSlNfTk9STUFMUgZqc3R5cGUSGQoEbGF6eRgFIAEoCDoFZmFsc2VSBGxhenkSLgoPdW52ZXJpZmllZF9sYXp5GA8gASgIOgVmYWxzZVIOdW52ZXJpZmllZExhenkSJQoKZGVwcmVjYXRlZBgDIAEoCDoFZmFsc2VSCmRlcHJlY2F0ZWQSGQoEd2VhaxgKIAEoCDoFZmFsc2VSBHdlYWsSWAoUdW5pbnRlcnByZXRlZF9vcHRpb24Y5wcgAygLMiQuZ29vZ2xlLnByb3RvYnVmLlVuaW50ZXJwcmV0ZWRPcHRpb25SE3VuaW50ZXJwcmV0ZWRPcHRpb24iLwoFQ1R5cGUSCgoGU1RSSU5HEAASCAoEQ09SRBABEhAKDFNUUklOR19QSUVDRRACIjUKBkpTVHlwZRINCglKU19OT1JNQUwQABINCglKU19TVFJJTkcQARINCglKU19OVU1CRVIQAioJCOgHEICAgIACSgQIBBAFInMKDE9uZW9mT3B0aW9ucxJYChR1bmludGVycHJldGVkX29wdGlvbhjnByADKAsyJC5nb29nbGUucHJvdG9idWYuVW5pbnRlcnByZXRlZE9wdGlvblITdW5pbnRlcnByZXRlZE9wdGlvbioJCOgHEICAgIACIsABCgtFbnVtT3B0aW9ucxIfCgthbGxvd19hbGlhcxgCIAEoCFIKYWxsb3dBbGlhcxIlCgpkZXByZWNhdGVkGAMgASgIOgVmYWxzZVIKZGVwcmVjYXRlZBJYChR1bmludGVycHJldGVkX29wdGlvbhjnByADKAsyJC5nb29nbGUucHJvdG9idWYuVW5pbnRlcnByZXRlZE9wdGlvblITdW5pbnRlcnByZXRlZE9wdGlvbioJCOgHEICAgIACSgQIBRAGIp4BChBFbnVtVmFsdWVPcHRpb25zEiUKCmRlcHJlY2F0ZWQYASABKAg6BWZhbHNlUgpkZXByZWNhdGVkElgKFHVuaW50ZXJwcmV0ZWRfb3B0aW9uGOcHIAMoCzIkLmdvb2dsZS5wcm90b2J1Zi5VbmludGVycHJldGVkT3B0aW9uUhN1bmludGVycHJldGVkT3B0aW9uKgkI6AcQgICAgAIinAEKDlNlcnZpY2VPcHRpb25zEiUKCmRlcHJlY2F0ZWQYISABKAg6BWZhbHNlUgpkZXByZWNhdGVkElgKFHVuaW50ZXJwcmV0ZWRfb3B0aW9uGOcHIAMoCzIkLmdvb2dsZS5wcm90b2J1Zi5VbmludGVycHJldGVkT3B0aW9uUhN1bmludGVycHJldGVkT3B0aW9uKgkI6AcQgICAgAIi4AIKDU1ldGhvZE9wdGlvbnMSJQoKZGVwcmVjYXRlZBghIAEoCDoFZmFsc2VSCmRlcHJlY2F0ZWQScQoRaWRlbXBvdGVuY3lfbGV2ZWwYIiABKA4yLy5nb29nbGUucHJvdG9idWYuTWV0aG9kT3B0aW9ucy5JZGVtcG90ZW5jeUxldmVsOhNJREVNUE9URU5DWV9VTktOT1dOUhBpZGVtcG90ZW5jeUxldmVsElgKFHVuaW50ZXJwcmV0ZWRfb3B0aW9uGOcHIAMoCzIkLmdvb2dsZS5wcm90b2J1Zi5VbmludGVycHJldGVkT3B0aW9uUhN1bmludGVycHJldGVkT3B0aW9uIlAKEElkZW1wb3RlbmN5TGV2ZWwSFwoTSURFTVBPVEVOQ1lfVU5LTk9XThAAEhMKD05PX1NJREVfRUZGRUNUUxABEg4KCklERU1QT1RFTlQQAioJCOgHEICAgIACIpoDChNVbmludGVycHJldGVkT3B0aW9uEkEKBG5hbWUYAiADKAsyLS5nb29nbGUucHJvdG9idWYuVW5pbnRlcnByZXRlZE9wdGlvbi5OYW1lUGFydFIEbmFtZRIpChBpZGVudGlmaWVyX3ZhbHVlGAMgASgJUg9pZGVudGlmaWVyVmFsdWUSLAoScG9zaXRpdmVfaW50X3ZhbHVlGAQgASgEUhBwb3NpdGl2ZUludFZhbHVlEiwKEm5lZ2F0aXZlX2ludF92YWx1ZRgFIAEoA1IQbmVnYXRpdmVJbnRWYWx1ZRIhCgxkb3VibGVfdmFsdWUYBiABKAFSC2RvdWJsZVZhbHVlEiEKDHN0cmluZ192YWx1ZRgHIAEoDFILc3RyaW5nVmFsdWUSJwoPYWdncmVnYXRlX3ZhbHVlGAggASgJUg5hZ2dyZWdhdGVWYWx1ZRpKCghOYW1lUGFydBIbCgluYW1lX3BhcnQYASACKAlSCG5hbWVQYXJ0EiEKDGlzX2V4dGVuc2lvbhgCIAIoCFILaXNFeHRlbnNpb24ipwIKDlNvdXJjZUNvZGVJbmZvEkQKCGxvY2F0aW9uGAEgAygLMiguZ29vZ2xlLnByb3RvYnVmLlNvdXJjZUNvZGVJbmZvLkxvY2F0aW9uUghsb2NhdGlvbhrOAQoITG9jYXRpb24SFgoEcGF0aBgBIAMoBUICEAFSBHBhdGgSFgoEc3BhbhgCIAMoBUICEAFSBHNwYW4SKQoQbGVhZGluZ19jb21tZW50cxgDIAEoCVIPbGVhZGluZ0NvbW1lbnRzEisKEXRyYWlsaW5nX2NvbW1lbnRzGAQgASgJUhB0cmFpbGluZ0NvbW1lbnRzEjoKGWxlYWRpbmdfZGV0YWNoZWRfY29tbWVudHMYBiADKAlSF2xlYWRpbmdEZXRhY2hlZENvbW1lbnRzItEBChFHZW5lcmF0ZWRDb2RlSW5mbxJNCgphbm5vdGF0aW9uGAEgAygLMi0uZ29vZ2xlLnByb3RvYnVmLkdlbmVyYXRlZENvZGVJbmZvLkFubm90YXRpb25SCmFubm90YXRpb24abQoKQW5ub3RhdGlvbhIWCgRwYXRoGAEgAygFQgIQAVIEcGF0aBIfCgtzb3VyY2VfZmlsZRgCIAEoCVIKc291cmNlRmlsZRIUCgViZWdpbhgDIAEoBVIFYmVnaW4SEAoDZW5kGAQgASgFUgNlbmRCfgoTY29tLmdvb2dsZS5wcm90b2J1ZkIQRGVzY3JpcHRvclByb3Rvc0gBWi1nb29nbGUuZ29sYW5nLm9yZy9wcm90b2J1Zi90eXBlcy9kZXNjcmlwdG9ycGL4AQGiAgNHUEKqAhpHb29nbGUuUHJvdG9idWYuUmVmbGVjdGlvbgrYAQocZ29vZ2xlL2FwaS9hbm5vdGF0aW9ucy5wcm90bxIKZ29vZ2xlLmFwaRoVZ29vZ2xlL2FwaS9odHRwLnByb3RvGiBnb29nbGUvcHJvdG9idWYvZGVzY3JpcHRvci5wcm90bzpLCgRodHRwEh4uZ29vZ2xlLnByb3RvYnVmLk1ldGhvZE9wdGlvbnMYsMq8IiABKAsyFC5nb29nbGUuYXBpLkh0dHBSdWxlUgRodHRwQh5aHHByb3RvL3RoaXJkX3BhcnR5L2dvb2dsZS9hcGliBnByb3RvMwrcAwoHYS5wcm90bxIFaGVsbG8aHGdvb2dsZS9hcGkvYW5ub3RhdGlvbnMucHJvdG8iHQoHUmVxdWVzdBISCgRwaW5nGAEgASgJUgRwaW5nIh4KCFJlc3BvbnNlEhIKBHBvbmcYASABKAlSBHBvbmcy2QIKBUhlbGxvEkQKB1BpbmdHZXQSDi5oZWxsby5SZXF1ZXN0Gg8uaGVsbG8uUmVzcG9uc2UiGILT5JMCEhINL3YxL2dldC97Zm9vfToBKhJACghQaW5nUG9zdBIOLmhlbGxvLlJlcXVlc3QaDy5oZWxsby5SZXNwb25zZSITgtPkkwINIggvdjEvcG9zdDoBKhI+CgdQaW5nUHV0Eg4uaGVsbG8uUmVxdWVzdBoPLmhlbGxvLlJlc3BvbnNlIhKC0+STAgwaBy92MS9wdXQ6ASoSRAoKUGluZ0RlbGV0ZRIOLmhlbGxvLlJlcXVlc3QaDy5oZWxsby5SZXNwb25zZSIVgtPkkwIPKgovdjEvZGVsZXRlOgEqEkIKCVBpbmdQYXRjaBIOLmhlbGxvLlJlcXVlc3QaDy5oZWxsby5SZXNwb25zZSIUgtPkkwIOMgkvdjEvcGF0Y2g6ASpCCVoHLi9oZWxsb2IGcHJvdG8z`
)

func TestGetMethods(t *testing.T) {
	tmpfile, err := os.CreateTemp(os.TempDir(), hash.Md5Hex([]byte(b64pb)))
	assert.Nil(t, err)
	b, err := base64.StdEncoding.DecodeString(b64pb)
	assert.Nil(t, err)
	assert.Nil(t, os.WriteFile(tmpfile.Name(), b, os.ModeTemporary))
	defer os.Remove(tmpfile.Name())

	source, err := grpcurl.DescriptorSourceFromProtoSets(tmpfile.Name())
	assert.Nil(t, err)
	methods, err := GetMethods(source)
	assert.Nil(t, err)
	assert.EqualValues(t, []Method{
		{
			RpcPath: "hello.Hello/Ping",
		},
	}, methods)
}

func TestGetMethodsWithAnnotations(t *testing.T) {
	tmpfile, err := os.CreateTemp(os.TempDir(), hash.Md5Hex([]byte(b64pb)))
	assert.Nil(t, err)
	b, err := base64.StdEncoding.DecodeString(b64pbWithAnnotations)
	assert.Nil(t, err)
	assert.Nil(t, os.WriteFile(tmpfile.Name(), b, os.ModeTemporary))
	defer os.Remove(tmpfile.Name())

	source, err := grpcurl.DescriptorSourceFromProtoSets(tmpfile.Name())
	assert.Nil(t, err)
	methods, err := GetMethods(source)
	assert.Nil(t, err)
	assert.EqualValues(t, []Method{
		{
			HttpMethod: http.MethodGet,
			HttpPath:   "/v1/get/:foo",
			RpcPath:    "hello.Hello/PingGet",
		},
		{
			HttpMethod: http.MethodPost,
			HttpPath:   "/v1/post",
			RpcPath:    "hello.Hello/PingPost",
		},
		{
			HttpMethod: http.MethodPut,
			HttpPath:   "/v1/put",
			RpcPath:    "hello.Hello/PingPut",
		},
		{
			HttpMethod: http.MethodDelete,
			HttpPath:   "/v1/delete",
			RpcPath:    "hello.Hello/PingDelete",
		},
		{
			HttpMethod: http.MethodPatch,
			HttpPath:   "/v1/patch",
			RpcPath:    "hello.Hello/PingPatch",
		},
	}, methods)
}

func TestGetMethodsBadCases(t *testing.T) {
	t.Run("no services", func(t *testing.T) {
		source := &mockDescriptorSource{
			servicesErr: errors.New("no services"),
		}
		_, err := GetMethods(source)
		assert.NotNil(t, err)
	})

	t.Run("no symbol in services", func(t *testing.T) {
		source := &mockDescriptorSource{
			services:  []string{"hello.Hello"},
			symbolErr: errors.New("no symbol"),
		}
		_, err := GetMethods(source)
		assert.NotNil(t, err)
	})

	t.Run("no symbol in services", func(t *testing.T) {
		source := &mockDescriptorSource{
			services:  []string{"hello.Hello"},
			symbolErr: errors.New("no symbol"),
		}
		_, err := GetMethods(source)
		assert.NotNil(t, err)
	})
}

type mockDescriptorSource struct {
	symbolDesc  desc.Descriptor
	symbolErr   error
	services    []string
	servicesErr error
}

func (m *mockDescriptorSource) AllExtensionsForType(_ string) ([]*desc.FieldDescriptor, error) {
	return nil, nil
}

func (m *mockDescriptorSource) FindSymbol(_ string) (desc.Descriptor, error) {
	return m.symbolDesc, m.symbolErr
}

func (m *mockDescriptorSource) ListServices() ([]string, error) {
	return m.services, m.servicesErr
}
