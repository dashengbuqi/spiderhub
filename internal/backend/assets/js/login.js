var delta = 0x9E3779B9;

this.longArrayToString = function(data, includeLength){
    var length = data.length;
    var n = (length - 1) << 2;
    if (includeLength) {
        var m = data[length - 1];
        if ((m < n - 3) || (m > n))
            return null;
        n = m;
    }
    for (var i = 0; i < length; i++) {
        data[i] = String.fromCharCode(data[i] & 0xff, data[i] >>> 8 & 0xff, data[i] >>> 16 & 0xff, data[i] >>> 24 & 0xff);
    }
    if (includeLength) {
        return data.join('').substring(0, n);
    }
    else {
        return data.join('');
    }
}

this.stringToLongArray = function(string, includeLength){
    var length = string.length;
    var result = [];
    for (var i = 0; i < length; i += 4) {
        result[i >> 2] = string.charCodeAt(i) |
            string.charCodeAt(i + 1) << 8 |
            string.charCodeAt(i + 2) << 16 |
            string.charCodeAt(i + 3) << 24;
    }
    if (includeLength) {
        result[result.length] = length;
    }
    return result;
};
this.encrypt = function(string, key){
    if (string == "") {
        return "";
    }
    var v = stringToLongArray(string, true);
    var k = stringToLongArray(key, false);
    if (k.length < 4) {
        k.length = 4;
    }
    var n = v.length - 1;

    var z = v[n], y = v[0];
    var mx, e, p, q = Math.floor(6 + 52 / (n + 1)), sum = 0;
    while (0 < q--) {
        sum = sum + delta & 0xffffffff;
        e = sum >>> 2 & 3;
        for (p = 0; p < n; p++) {
            y = v[p + 1];
            mx = (z >>> 5 ^ y << 2) + (y >>> 3 ^ z << 4) ^ (sum ^ y) + (k[p & 3 ^ e] ^ z);
            z = v[p] = v[p] + mx & 0xffffffff;
        }
        y = v[0];
        mx = (z >>> 5 ^ y << 2) + (y >>> 3 ^ z << 4) ^ (sum ^ y) + (k[p & 3 ^ e] ^ z);
        z = v[n] = v[n] + mx & 0xffffffff;
    }

    return longArrayToString(v, false);
}

this.decrypt = function(string, key){
    if (string == "") {
        return "";
    }
    var v = stringToLongArray(string, false);
    var k = stringToLongArray(key, false);
    if (k.length < 4) {
        k.length = 4;
    }
    var n = v.length - 1;

    var z = v[n - 1], y = v[0];
    var mx, e, p, q = Math.floor(6 + 52 / (n + 1)), sum = q * delta & 0xffffffff;
    while (sum != 0) {
        e = sum >>> 2 & 3;
        for (p = n; p > 0; p--) {
            z = v[p - 1];
            mx = (z >>> 5 ^ y << 2) + (y >>> 3 ^ z << 4) ^ (sum ^ y) + (k[p & 3 ^ e] ^ z);
            y = v[p] = v[p] - mx & 0xffffffff;
        }
        z = v[n];
        mx = (z >>> 5 ^ y << 2) + (y >>> 3 ^ z << 4) ^ (sum ^ y) + (k[p & 3 ^ e] ^ z);
        y = v[0] = v[0] - mx & 0xffffffff;
        sum = sum - delta & 0xffffffff;
    }

    return longArrayToString(v, true);
}

if (typeof(base64Encode) == "undefined") {
    base64Encode = function(){
        var base64EncodeChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'.split('');
        return function(str){
            var out, i, j, len, r, l, c;
            i = j = 0;
            len = str.length;
            r = len % 3;
            len = len - r;
            l = (len / 3) << 2;
            if (r > 0) {
                l += 4;
            }
            out = new Array(l);

            while (i < len) {
                c = str.charCodeAt(i++) << 16 |
                    str.charCodeAt(i++) << 8 |
                    str.charCodeAt(i++);
                out[j++] = base64EncodeChars[c >> 18] +
                    base64EncodeChars[c >> 12 & 0x3f] +
                    base64EncodeChars[c >> 6 & 0x3f] +
                    base64EncodeChars[c & 0x3f];
            }
            if (r == 1) {
                c = str.charCodeAt(i++);
                out[j++] = base64EncodeChars[c >> 2] +
                    base64EncodeChars[(c & 0x03) << 4] +
                    "==";
            }
            else
            if (r == 2) {
                c = str.charCodeAt(i++) << 8 |
                    str.charCodeAt(i++);
                out[j++] = base64EncodeChars[c >> 10] +
                    base64EncodeChars[c >> 4 & 0x3f] +
                    base64EncodeChars[(c & 0x0f) << 2] +
                    "=";
            }
            return out.join('');
        }
    }();
}

if (typeof(base64Decode) == "undefined") {
    base64Decode = function(){
        var base64DecodeChars = [-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, -1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1];
        return function(str){
            var c1, c2, c3, c4;
            var i, j, len, r, l, out;

            len = str.length;
            if (len % 4 != 0) {
                return '';
            }
            if (/[^ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789\+\/\=]/.test(str)) {
                return '';
            }
            if (str.charAt(len - 2) == '=') {
                r = 1;
            }
            else
            if (str.charAt(len - 1) == '=') {
                r = 2;
            }
            else {
                r = 0;
            }
            l = len;
            if (r > 0) {
                l -= 4;
            }
            l = (l >> 2) * 3 + r;
            out = new Array(l);

            i = j = 0;
            while (i < len) {
                // c1
                c1 = base64DecodeChars[str.charCodeAt(i++)];
                if (c1 == -1)
                    break;

                // c2
                c2 = base64DecodeChars[str.charCodeAt(i++)];
                if (c2 == -1)
                    break;

                out[j++] = String.fromCharCode((c1 << 2) | ((c2 & 0x30) >> 4));

                // c3
                c3 = base64DecodeChars[str.charCodeAt(i++)];
                if (c3 == -1)
                    break;

                out[j++] = String.fromCharCode(((c2 & 0x0f) << 4) | ((c3 & 0x3c) >> 2));

                // c4
                c4 = base64DecodeChars[str.charCodeAt(i++)];
                if (c4 == -1)
                    break;

                out[j++] = String.fromCharCode(((c3 & 0x03) << 6) | c4);
            }
            return out.join('');
        }
    }();
}
String.prototype.startsWith = function(prefix){
    return this.substring(0, prefix.length) == prefix;
}
String.prototype.endsWith = function(suffix){
    return this.substring(this.length - suffix.length) == suffix;
}

String.prototype.toCharArray = function(){
    var charArr = new Array();
    for (var i = 0; i < this.length; i++)
        charArr[i] = this.charAt(i);
    return charArr;
}
String.prototype.trim = function(){
    return this.replace(/(^\s*)|(\s*$)/g, "");
}
String.prototype.replaceAll = function(AFindText, ARepText){
    raRegExp = new RegExp(AFindText, "g");
    return this.replace(raRegExp, ARepText);
}
String.prototype.ltrim = function(){
    return this.replace(/(^\s*)/g, "");
}
String.prototype.rtrim = function(){
    return this.replace(/(\s*$)/g, "");
}
String.prototype.base64Encode = function(){
    return base64Encode(this);
}
String.prototype.base64Decode = function(){
    return Base64Decode(this);
}
String.prototype.encrypt = function(key){
    return encodeURIComponent(base64Encode(encrypt(this, key)));
}
String.prototype.decrypt = function(key){
    return decrypt(base64Decode(decodeURIComponent(this)), key);
}
String.prototype.toUTF8 = function(){
    var str = this;
    if (str.match(/^[\x00-\x7f]*$/) != null) {
        return str.toString();
    }
    var out, i, j, len, c, c2;
    out = [];
    len = str.length;
    for (i = 0, j = 0; i < len; i++, j++) {
        c = str.charCodeAt(i);
        if (c <= 0x7f) {
            out[j] = str.charAt(i);
        }
        else
        if (c <= 0x7ff) {
            out[j] = String.fromCharCode(0xc0 | (c >>> 6), 0x80 | (c & 0x3f));
        }
        else
        if (c < 0xd800 || c > 0xdfff) {
            out[j] = String.fromCharCode(0xe0 | (c >>> 12), 0x80 | ((c >>> 6) & 0x3f), 0x80 | (c & 0x3f));
        }
        else {
            if (++i < len) {
                c2 = str.charCodeAt(i);
                if (c <= 0xdbff && 0xdc00 <= c2 && c2 <= 0xdfff) {
                    c = ((c & 0x03ff) << 10 | (c2 & 0x03ff)) + 0x010000;
                    if (0x010000 <= c && c <= 0x10ffff) {
                        out[j] = String.fromCharCode(0xf0 | ((c >>> 18) & 0x3f), 0x80 | ((c >>> 12) & 0x3f), 0x80 | ((c >>> 6) & 0x3f), 0x80 | (c & 0x3f));
                    }
                    else {
                        out[j] = '?';
                    }
                }
                else {
                    i--;
                    out[j] = '?';
                }
            }
            else {
                i--;
                out[j] = '?';
            }
        }
    }
    return out.join('');
}

String.prototype.sha256 = function(){
    var rotateRight = function(n, x){
        return ((x >>> n) | (x << (32 - n)));
    }
    var choice = function(x, y, z){
        return ((x & y) ^ (~ x & z));
    }
    var majority = function(x, y, z){
        return ((x & y) ^ (x & z) ^ (y & z));
    }
    var sha256_Sigma0 = function(x){
        return (rotateRight(2, x) ^ rotateRight(13, x) ^ rotateRight(22, x));
    }
    var sha256_Sigma1 = function(x){
        return (rotateRight(6, x) ^ rotateRight(11, x) ^ rotateRight(25, x));
    }
    var sha256_sigma0 = function(x){
        return (rotateRight(7, x) ^ rotateRight(18, x) ^ (x >>> 3));
    }
    var sha256_sigma1 = function(x){
        return (rotateRight(17, x) ^ rotateRight(19, x) ^ (x >>> 10));
    }
    var sha256_expand = function(W, j){
        return (W[j & 0x0f] += sha256_sigma1(W[(j + 14) & 0x0f]) + W[(j + 9) & 0x0f] +
            sha256_sigma0(W[(j + 1) & 0x0f]));
    }

    /* Hash constant words K: */
    var K256 = new Array(0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5, 0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174, 0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da, 0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967, 0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070, 0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3, 0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2);

    /* global arrays */
    var ihash, count, buffer;
    var sha256_hex_digits = "0123456789abcdef";

    /* Add 32-bit integers with 16-bit operations (bug in some JS-interpreters: 
     overflow) */
    var safe_add = function(x, y){
        var lsw = (x & 0xffff) + (y & 0xffff);
        var msw = (x >> 16) + (y >> 16) + (lsw >> 16);
        return (msw << 16) | (lsw & 0xffff);
    }

    /* Initialise the SHA256 computation */
    var sha256_init = function(){
        ihash = new Array(8);
        count = new Array(2);
        buffer = new Array(64);
        count[0] = count[1] = 0;
        ihash[0] = 0x6a09e667;
        ihash[1] = 0xbb67ae85;
        ihash[2] = 0x3c6ef372;
        ihash[3] = 0xa54ff53a;
        ihash[4] = 0x510e527f;
        ihash[5] = 0x9b05688c;
        ihash[6] = 0x1f83d9ab;
        ihash[7] = 0x5be0cd19;
    }

    /* Transform a 512-bit message block */
    var sha256_transform = function(){
        var a, b, c, d, e, f, g, h, T1, T2;
        var W = new Array(16);

        /* Initialize registers with the previous intermediate value */
        a = ihash[0];
        b = ihash[1];
        c = ihash[2];
        d = ihash[3];
        e = ihash[4];
        f = ihash[5];
        g = ihash[6];
        h = ihash[7];

        /* make 32-bit words */
        for (var i = 0; i < 16; i++)
            W[i] = ((buffer[(i << 2) + 3]) | (buffer[(i << 2) + 2] << 8) |
            (buffer[(i << 2) + 1] <<
            16) |
            (buffer[i << 2] << 24));

        for (var j = 0; j < 64; j++) {
            T1 = h + sha256_Sigma1(e) + choice(e, f, g) + K256[j];
            if (j < 16)
                T1 += W[j];
            else
                T1 += sha256_expand(W, j);
            T2 = sha256_Sigma0(a) + majority(a, b, c);
            h = g;
            g = f;
            f = e;
            e = safe_add(d, T1);
            d = c;
            c = b;
            b = a;
            a = safe_add(T1, T2);
        }

        /* Compute the current intermediate hash value */
        ihash[0] += a;
        ihash[1] += b;
        ihash[2] += c;
        ihash[3] += d;
        ihash[4] += e;
        ihash[5] += f;
        ihash[6] += g;
        ihash[7] += h;
    }

    /* Read the next chunk of data and update the SHA256 computation */
    var sha256_update = function(data, inputLen){
        var i, index, curpos = 0;
        /* Compute number of bytes mod 64 */
        index = ((count[0] >> 3) & 0x3f);
        var remainder = (inputLen & 0x3f);

        /* Update number of bits */
        if ((count[0] += (inputLen << 3)) < (inputLen << 3))
            count[1]++;
        count[1] += (inputLen >> 29);

        /* Transform as many times as possible */
        for (i = 0; i + 63 < inputLen; i += 64) {
            for (var j = index; j < 64; j++)
                buffer[j] = data.charCodeAt(curpos++);
            sha256_transform();
            index = 0;
        }

        /* Buffer remaining input */
        for (var j = 0; j < remainder; j++)
            buffer[j] = data.charCodeAt(curpos++);
    }

    /* Finish the computation by operations such as padding */
    var sha256_final = function(){
        var index = ((count[0] >> 3) & 0x3f);
        buffer[index++] = 0x80;
        if (index <= 56) {
            for (var i = index; i < 56; i++)
                buffer[i] = 0;
        }
        else {
            for (var i = index; i < 64; i++)
                buffer[i] = 0;
            sha256_transform();
            for (var i = 0; i < 56; i++)
                buffer[i] = 0;
        }
        buffer[56] = (count[1] >>> 24) & 0xff;
        buffer[57] = (count[1] >>> 16) & 0xff;
        buffer[58] = (count[1] >>> 8) & 0xff;
        buffer[59] = count[1] & 0xff;
        buffer[60] = (count[0] >>> 24) & 0xff;
        buffer[61] = (count[0] >>> 16) & 0xff;
        buffer[62] = (count[0] >>> 8) & 0xff;
        buffer[63] = count[0] & 0xff;
        sha256_transform();
    }

    /* Split the internal hash values into an array of bytes */
    var sha256_encode_bytes = function(){
        var j = 0;
        var output = new Array(32);
        for (var i = 0; i < 8; i++) {
            output[j++] = ((ihash[i] >>> 24) & 0xff);
            output[j++] = ((ihash[i] >>> 16) & 0xff);
            output[j++] = ((ihash[i] >>> 8) & 0xff);
            output[j++] = (ihash[i] & 0xff);
        }
        return output;
    }

    /* Get the internal hash as a hex string */
    var sha256_encode_hex = function(){
        var output = new String();
        for (var i = 0; i < 8; i++) {
            for (var j = 28; j >= 0; j -= 4)
                output += sha256_hex_digits.charAt((ihash[i] >>> j) & 0x0f);
        }
        return output;
    }

    /* Main function: returns a hex string representing the SHA256 value of the 
     given data */
    var sha256_digest = function(data){
        sha256_init();
        sha256_update(data, data.length);
        sha256_final();
        return sha256_encode_hex();
    }

    /* test if the JS-interpreter is working properly */
    var sha256_self_test = function(){
        return sha256_digest("message digest") ==
            "f7846f55cf23e14eebeab5b4e1550cad5b509e3348fbc4efa3a1413d393cb650";
    }

    return sha256_digest(this);
}
