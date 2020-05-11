import os
import importlib

class cloud_storage(object):
    def setObj(self, key: str, data: bytes):
        raise NotImplementedError
    def getObj(self, key: str) -> bytes:
        raise NotImplementedError 
    def delObj(self, key: str):
        raise NotImplementedError

class oss_storage(cloud_storage):
    def __init__(self, accessKeyId, accessKeySecret, securityToken, region):
        super().__init__()
        oss2 = importlib.import_module('oss2')
        auth = oss2.StsAuth(accessKeyId, accessKeySecret, securityToken)
        self.oss_client = oss2.Bucket(auth, 'oss-%s-internal.aliyuncs.com' % region, 'hcloudstorage')
    def setObj(self, key: str, data: bytes):
        self.oss_client.put_object(key, data)
    def getObj(self, key: str) -> bytes:
        return self.oss_client.get_object(key).read()
    def delObj(self, key: str):
        self.oss_client.delete_object(key)

def NewCloudStorage() -> cloud_storage:
    envs = os.environ
    env = envs.get('jointfaas_env', '')
    if env == 'aliyun':
        return oss_storage(envs['accessKeyId'], envs['accessKeySecret'], envs['securityToken'], envs['region'])
    elif env == 'aws':
        pass
    elif env == 'private':
        pass