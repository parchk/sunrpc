struct intpair {
	int a;
	int b;
};

const NFS3_FHSIZE    = 64;    /* Maximum bytes in a V3 file handle */
const NFS3_WRITEVERFSIZE = 8;
const NFS3_CREATEVERFSIZE = 8;
const NFS3_COOKIEVERFSIZE = 8;

typedef opaque cookieverf3[NFS3_COOKIEVERFSIZE];

typedef uint64_t cookie3;

struct nfs_fh3 {
    opaque       data<NFS3_FHSIZE>;
};

typedef string filename3<>;

struct READDIRPLUS3args {
    nfs_fh3      dir;
    cookie3      cookie;
    cookieverf3  cookieverf;
    count3       dircount;
    count3       maxcount;
};

struct entryplus3 {
    fileid3      fileid;
    filename3    name;
    cookie3      cookie;
    post_op_attr name_attributes;
    post_op_fh3  name_handle;
    entryplus3   *nextentry;
};

struct dirlistplus3 {
    entryplus3   *entries;
    bool         eof; 
};

struct READDIRPLUS3resok {
    post_op_attr dir_attributes;
    cookieverf3  cookieverf;
    dirlistplus3 reply;
};


struct READDIRPLUS3resfail {
    post_op_attr dir_attributes;
};

union READDIRPLUS3res switch (nfsstat3 status) {
    case NFS3_OK:
        READDIRPLUS3resok   resok;
    default:
        READDIRPLUS3resfail resfail;
};

enum ftype3 {
    NF3REG    = 1,
    NF3DIR    = 2,
    NF3BLK    = 3,
    NF3CHR    = 4,
    NF3LNK    = 5,
    NF3SOCK   = 6,
    NF3FIFO   = 7
};

typedef unsigned int mode3;

typedef unsigned int uid3;

typedef unsigned int gid3;

typedef uint64_t size3;

typedef uint64_t fileid3;

struct specdata3 {
    unsigned int specdata1;
    unsigned int specdata2;
};

struct nfstime3 {
    unsigned int seconds;
    unsigned int nseconds;
};

struct fattr3 {
    ftype3       type;
    mode3        mode;
    unsigned int nlink;
    uid3         uid;
    gid3         gid;
    size3        size;
    size3        used;
    specdata3    rdev;
    uint64_t     fsid;
    fileid3      fileid;
    nfstime3     atime;
    nfstime3     mtime;
    nfstime3     ctime;
};

union post_op_attr switch (bool attributes_follow) {
    case TRUE:
        fattr3   attributes;
    case FALSE:
        void;
};


enum nfsstat3 {
    NFS3_OK             = 0,
    NFS3ERR_PERM        = 1,
    NFS3ERR_NOENT       = 2,
    NFS3ERR_IO          = 5,
    NFS3ERR_NXIO        = 6,
    NFS3ERR_ACCES       = 13,
    NFS3ERR_EXIST       = 17,
    NFS3ERR_XDEV        = 18,
    NFS3ERR_NODEV       = 19,
    NFS3ERR_NOTDIR      = 20,
    NFS3ERR_ISDIR       = 21,
    NFS3ERR_INVAL       = 22,
    NFS3ERR_FBIG        = 27,
    NFS3ERR_NOSPC       = 28,
    NFS3ERR_ROFS        = 30,
    NFS3ERR_MLINK       = 31,
    NFS3ERR_NAMETOOLONG = 63,
    NFS3ERR_NOTEMPTY    = 66,
    NFS3ERR_DQUOT       = 69,
    NFS3ERR_STALE       = 70,
    NFS3ERR_REMOTE      = 71,
    NFS3ERR_BADHANDLE   = 10001,
    NFS3ERR_NOT_SYNC    = 10002,
    NFS3ERR_BAD_COOKIE  = 10003,
    NFS3ERR_NOTSUPP     = 10004,
    NFS3ERR_TOOSMALL    = 10005,
    NFS3ERR_SERVERFAULT = 10006,
    NFS3ERR_BADTYPE     = 10007,
    NFS3ERR_JUKEBOX     = 10008
};

typedef uint64_t offset3;

typedef unsigned int count3;


program ARITH_PROG {
	version ARITH_VERS {
		int ADD(intpair) = 1;
		int MULTIPLY(intpair) = 2;
		READDIRPLUS3res TEST_READ(READDIRPLUS3args) =3;
	} = 1;
} = 12345;
