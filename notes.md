All communications in the DNS protocol are carried in a single format called a "message". Each message consists of 5 sections: header, question, answer, authority, and an additional space.

Packet Identifier (ID)	16 bits	A random ID assigned to query packets. Response packets must reply with the same ID.
Expected value: 1234.
Query/Response Indicator (QR)	1 bit	1 for a reply packet, 0 for a question packet.
Expected value: 1.
Operation Code (OPCODE)	4 bits	Specifies the kind of query in a message.
Expected value: 0.
Authoritative Answer (AA)	1 bit	1 if the responding server "owns" the domain queried, i.e., it's authoritative.
Expected value: 0.
Truncation (TC)	1 bit	1 if the message is larger than 512 bytes. Always 0 in UDP responses.
Expected value: 0.
Recursion Desired (RD)	1 bit	Sender sets this to 1 if the server should recursively resolve this query, 0 otherwise.
Expected value: 0.
Recursion Available (RA)	1 bit	Server sets this to 1 to indicate that recursion is available.
Expected value: 0.
Reserved (Z)	3 bits	Used by DNSSEC queries. At inception, it was reserved for future use.
Expected value: 0.
Response Code (RCODE)	4 bits	Response code indicating the status of the response.
Expected value: 0 (no error).
Question Count (QDCOUNT)	16 bits	Number of questions in the Question section.
Expected value: 0.
Answer Record Count (ANCOUNT)	16 bits	Number of records in the Answer section.
Expected value: 0.
Authority Record Count (NSCOUNT)	16 bits	Number of records in the Authority section.
Expected value: 0.
Additional Record Count (ARCOUNT)	16 bits	Number of records in the Additional section.
Expected value: 0.
The header section is always 12 bytes long. Integers are encoded in big-endian format.