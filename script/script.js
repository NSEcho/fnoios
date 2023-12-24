// constants for file 
const O_RDWR = 2;

// functions
const open = getFunc("open", "int", ["pointer", "int"]);
const read = getFunc("read", "int", ["int", "pointer", "int"]);
const close = getFunc("close", "int", ["int"]);
const dup = getFunc("dup", "int", ["int"]);
const posix_openpt = getFunc("posix_openpt", "int",  ["int"]);
const grantpt = getFunc("grantpt", "int", ["int"]);
const unlockpt = getFunc("unlockpt", "int", ["int"]);
const ptsname = getFunc("ptsname", "pointer", ["int"]);

// create new native function
function getFunc(name, ret_type, args) {
    return new NativeFunction(Module.findExportByName(null, name), ret_type, args);
}

// JavaScript string to C string
function jstr(cs) {
    return Memory.readUtf8String(cs);
}

let ptyFD = null;
let ptt = null;
let psname = null;

function create_pty() {
    ptyFD = posix_openpt(O_RDWR);
    grantpt(ptyFD);
    unlockpt(ptyFD);
    psname = ptsname(ptyFD);
    ptt = open(psname, O_RDWR);
}

function run() {
    if (psname == null) {
        create_pty();
        send(JSON.stringify({
            "type": "info",
            "data": `Created pseudo-terminal ${jstr(psname)}\n`
        }));
        close(1);
        close(2);
        dup(ptyFD)
        dup(ptyFD);
    }
}

rpc.exports = {
    start() {
        run();
    },
    read() {
        let r = 0;
        let buffer = Memory.alloc(1);

        while ((r = read(ptt, buffer, 1)) > 0) {
            send(JSON.stringify({
                "type": "read",
                "data": jstr(buffer)
            }))
        }
    }
};
