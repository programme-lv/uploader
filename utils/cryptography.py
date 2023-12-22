import bcrypt
import hashlib


def bcrypt_password(password):
    salt = bcrypt.gensalt()
    hashed = bcrypt.hashpw(password.encode('utf-8'), salt)
    return hashed.decode('utf-8')


def sha256_string(data):
    sha256_hash = hashlib.sha256(data.encode('utf-8'))
    return sha256_hash.hexdigest()


def sha256_bytes(data):
    sha256_hash = hashlib.sha256(data)
    return sha256_hash.hexdigest()


def sha256_file(filename):
    sha256_hash = hashlib.sha256()
    with open(filename, "rb") as f:
        # Read and update hash string value in blocks of 4K
        for byte_block in iter(lambda: f.read(4096), b""):
            sha256_hash.update(byte_block)
        return sha256_hash.hexdigest()
