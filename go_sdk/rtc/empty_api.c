/**
 * empty_api.c
 * Empty API implementation
 * Function: empty API implementation, to avoid the dependency of the sdk
 * when the sdk is not used, or the sdk is not available, the empty API implementation can be used to avoid the compilation error
 * the empty API implementation is a placeholder, and the actual implementation is provided by the sdk
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>

#include "empty_api.h"


void agora_local_audio_track_set_send_delay_ms(void *cTrack, int delay_ms) {
    (void)cTrack;
    (void)delay_ms;
    return ;
}
int agora_local_audio_track_set_total_extra_send_ms(void *cTrack, uint64_t total_extra_send_ms) {
    (void)cTrack;
    (void)total_extra_send_ms;
    return 0;
}

