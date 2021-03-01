CREATE TABLE IF NOT EXISTS number (
	id INTEGER PRIMARY KEY,
	cc NCHAR(3) NOT NULL,
	ndc NCHAR(4) NOT NULL,
    sn NCHAR(13) NOT NULL, 
	used BOOLEAN NOT NULL DEFAULT 0, 
   	domain TEXT NOT NULL,
	carrier TEXT NOT NULL,
	userID  INTEGER NOT NULL DEFAULT 0, 
	allocated INTEGER NOT NULL DEFAULT 0, 
	reserved  INTEGER NOT NULL DEFAULT 0, 
	deallocated INTEGER NOT NULL DEFAULT 0, 
	portedIn  INTEGER NOT NULL DEFAULT 0, 
	portedOut INTEGER NOT NULL DEFAULT 0, 
    CONSTRAINT unq UNIQUE (cc, ndc, sn)
);

INSERT into number (id, cc, ndc, sn, domain, carrier) values (1, "353" , "086", "0111111", "test.com","anycarrier");
INSERT into number (id, cc, ndc, sn, domain, carrier, used, userID, allocated) values (2, "353" , "086", "0111112", "test.com","anycarrier", 1, 24, 1612564816);

