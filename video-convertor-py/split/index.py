# -*- coding: utf-8 -*-
import subprocess
from jfstorage import cloudstorage
import logging
import json
import os
import time
import math

LOGGER = logging.getLogger()

MAX_SPLIT_NUM = 100

ROOT = "tmp"
FFMPEG_BIN = "/tmp/ffmpeg"
FFPROBE_BIN = "/tmp/ffprobe"

class FFmpegError(Exception):
    def __init__(self, message, status):
        super().__init__(message, status)
        self.message = message
        self.status = status

def exec_FFmpeg_cmd(cmd_lst):
    try:
        subprocess.check_call(cmd_lst)
    except subprocess.CalledProcessError as exc:
        LOGGER.error('returncode:{}'.format(exc.returncode))
        LOGGER.error('cmd:{}'.format(exc.cmd))
        LOGGER.error('output:{}'.format(exc.output))
        # log json to Log Service as db
        # or insert record in mysql, etc
        raise FFmpegError(exc.output, exc.returncode)

def getVideoDuration(input_video):
    cmd = '{0} -i {1} -show_entries format=duration -v quiet -of csv="p=0"'.format(
        FFPROBE_BIN, input_video)
    raw_result = subprocess.check_output(cmd, shell=True)
    result = raw_result.decode().replace("\n", "").strip()
    duration = float(result)
    return duration

def downloadFFmpeg(cs):
    if os.path.exists(FFMPEG_BIN):
        return
    ffmpeg = open(FFMPEG_BIN, 'wb')
    ffmpeg.write(cs.getObj('ffmpeg'))
    ffmpeg.close()
    os.system("chmod 777 " + FFMPEG_BIN)
    ffprobe = open(FFPROBE_BIN, 'wb')
    ffprobe.write(cs.getObj('ffprobe'))
    ffprobe.close()
    os.system("chmod 777 " + FFPROBE_BIN)

def handler(event):
    cs = cloudstorage.NewCloudStorage()
    downloadFFmpeg(cs)
    video_key = event['video_key']
    segment_time_seconds = event['segment_time_seconds']
    
    video_path = ("/tmp/video.avi")
    video = open(video_path, "wb")
    video.write(cs.getObj(video_key))
    
    video_duration = getVideoDuration(video_path)
    split_num = math.ceil(video_duration/segment_time_seconds)
    if split_num > MAX_SPLIT_NUM:
        segment_time_seconds = int(math.ceil(video_duration/MAX_SPLIT_NUM)) + 1
    
    segment_time_seconds = str(segment_time_seconds)
    exec_FFmpeg_cmd([FFMPEG_BIN, '-i', video_path, "-c", "copy", "-f", "segment", "-segment_time",
                     segment_time_seconds, "-reset_timestamps", "1", ROOT + "/split_piece_" + video_key + "%02d.avi"])

    split_keys = []
    for filename in os.listdir(ROOT):
        if filename.startswith('split_'):
            split_keys.append(filename)
            f = open(os.path.join(ROOT, filename), 'rb')
            cs.setObj(filename, f.read())

    return {"split_keys": split_keys}