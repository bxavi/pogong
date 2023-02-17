import re
path = r"./sqlc";
inputfile = r"./sql/schema.sql"

outputfile = r"/query.sql"

#STORAGE FOR PRIMARY KEYS
keys = {}
sequences = {}

def makequeries(tablename, cols, path):
    with open(path, 'a', encoding='latin-1') as file:
        out = ""
        
        #SELECT STATEMENT
        if len(keys[tablename]) > 0: # only create select statements for tables with primary keys
            if tablename in keys and len(keys[tablename]) > 0: # only do a select statement if a primary key found 
                #SELECT STATEMENT
                out += "-- name: Get" + tablename.title() + " :one\n"
                out += "SELECT * FROM " + tablename + "\n"
                out += "WHERE " 
                #print(tablename.title().lower())
                
                for pkey in range( len(keys[tablename.title().lower()] )):
                    #print(keys[tablename.title().lower()][pkey])
                    out += keys[tablename.title().lower()][pkey] + " = $" + str(pkey + 1) + " AND "
                out = out[:-4]
                out += " LIMIT 1;\n\n"

        if "tag" in cols:
            #SELECT TAG STATEMENT
                out += "-- name: GetWithTag" + tablename.title() + " :one\n"
                out += "SELECT * FROM " + tablename + "\n"
                out += "WHERE tag = $1"
                out += " LIMIT 1;\n\n"

        #SELECT ALL STATEMENT
        out += "-- name: List" + tablename.title() + " :many\n"
        out += "SELECT * FROM " +  tablename + "\n"
        out += "ORDER BY " + cols[1] + ";\n\n"

        #INSERT STATEMENT
        out += "-- name: Create" + tablename.title() + " :one\n"
        out += "INSERT INTO " + tablename + " (\n"
        out += "\t"
        colscount = 0
        for x in range(0, len(cols)):
            if cols[x] not in sequences[tablename]: #if column is not a sequenced primary key
                #print(cols[x], "not in" + tablename + "  sequences ")
                out += cols[x]
                colscount += 1
                out += ", "
        out = out[:-2]
        out += "\n) VALUES (\n\t"
        for x in range (0, colscount):
            out += "$" + str(x+1) + ", "
        out = out[:-2]
        out += "\n)\nRETURNING *;\n\n"
        
        #UPDATE STATEMENT
        if len(keys[tablename]) == 1:
            skipUpdate = True
            for x in range(0, len(cols)):
                if cols[x] not in keys[tablename]: # if there is a column which are not keys (joins) then do update query 
                    skipUpdate = False
            
            if not skipUpdate:
                paracount = len(keys[tablename])# - len(sequences[tablename])
                out += "-- name: Update" + tablename.title() + " :one\n"
                out += "UPDATE " + tablename + "\n"
                out += "set "
                for x in range(0, len(cols)):
                    if cols[x] not in sequences[tablename] and cols[x] not in keys[tablename]:# and cols[x] not in keys[tablename]:
                        out += cols[x] + " = $" + str(paracount + len(sequences[tablename])) + ",\n"
                        paracount += 1
                out = out[:-2]
                out += "\nWHERE "
                
                for x in range(0, len(keys[tablename])):
                    out += keys[tablename][x] + " = $" + str(x+1) + " AND "

                out = out[:-5]
                out += "\nRETURNING *;\n\n"
                out += "\n"
        elif len(keys[tablename]) == 2:            
            paracount = len(keys[tablename])# - len(sequences[tablename])
            out += "-- name: Update" + tablename.title() + " :one\n"
            out += "UPDATE " + tablename + "\n"
            out += "set "
            for x in range(0, len(cols)):
                out += cols[x] + " = $" + str(x+1) + ",\n"
                paracount += 1
            out = out[:-2]
            out += "\nWHERE "
            
            for x in range(0, len(cols)):
                out += cols[x] + " = $" + str(x+1) + " AND "

            out = out[:-5]
            out += "\nRETURNING *;\n\n"
            out += "\n"
        if len(keys[tablename]) == 1 and "tag" in cols:
            skipUpdate = True
            for x in range(0, len(cols)):
                if cols[x] not in keys[tablename]: # if there is a column which are not keys (joins) then do update query 
                    skipUpdate = False
            
            if not skipUpdate:
                paracount = len(keys[tablename])# - len(sequences[tablename])
                out += "-- name: UpdateWithTag" + tablename.title() + " :one\n"
                out += "UPDATE " + tablename + "\n"
                out += "set "
                for x in range(0, len(cols)):
                    if cols[x] not in sequences[tablename] and cols[x] not in keys[tablename] and cols[x] != "tag":# and cols[x] not in keys[tablename]:
                        out += cols[x] + " = $" + str(paracount + len(sequences[tablename])) + ",\n"
                        paracount += 1
                out = out[:-2]
                out += "\nWHERE tag = $1\n"
                out += "\nRETURNING *;\n\n"
                out += "\n"


        #DELETE STATEMENT
        if len(keys[tablename]) == 1:
            out += "-- name: Delete" + tablename.title() + " :exec\n"
            out += "DELETE FROM " + tablename + "\n"
            out += "WHERE "
            #print(tablename.title().lower())
            for pkey in range( len(keys[tablename.title().lower()] )):
                #print(keys[tablename.title().lower()][pkey])
                out += keys[tablename.title().lower()][pkey] + " = $" + str(pkey + 1) + " AND "
            out = out[:-4]
            out += ";\n\n"
        
        if "tag" in cols:
            out += "-- name: DeleteWithTag" + tablename.title() + " :exec\n"
            out += "DELETE FROM " + tablename + "\n"
            out += "WHERE tag = $1;\n\n"

        elif len(keys[tablename]) == 2:
            out += "-- name: Delete" + tablename.title() + " :exec\n"
            out += "DELETE FROM " + tablename + "\n"
            out += "WHERE "
            #print(tablename.title().lower())

            for c in range( len(cols)):
                out += cols[c] + " = $" + str(c+1) + " AND "
            out = out[:-4]
            out += ";\n\n"

        file.write(out)

    return out

    #need next line to complete

def additionalQueries(path):
    with open(path, 'a', encoding='latin-1') as file:
        
        out = "-- Add queries here:\n"

        
        file.write(out)

f = open(inputfile,"r", encoding='latin-1')
lines = f.readlines()


#FIND PRIMARY KEYS
#pre pass to get primary keys
#NOTE THE NAME OF THE PRIMARY KEY CONSTRAINT MUST MATCH THE TABLE NAME I.E. USERS_PKEY FROM USERS TABLE.
for x in range(len(lines)):

    if len(lines[x]) > 0 and len(lines[x].split()) > 3 and lines[x].split()[0] == "ADD" and lines[x].split()[3] == "PRIMARY":
        table = lines[x].split()[2][:-5]
        #print("Table: ", table)
        keys[table] = []
        if table not in sequences:
            sequences[table] = []
        if table not in keys:
            keys[table] = []
        #print(lines[x].split())
        #print(lines[x].split('(')[1].split(')')[0])
        pkeys = lines[x].split('(')[1].split(')')[0].split()
        #print("Keys: ");
        for x in range(len(pkeys)):
            pkeys[x] = pkeys[x].replace(',','')
            #print(pkeys[x] )
            keys[table].append(pkeys[x])
            
    if len(lines[x]) > 0 and len(lines[x].split()) > 3 and lines[x].split()[0] == "ALTER" and lines[x].split()[1] == "SEQUENCE":
        #print("\n\nNEW Sequence detected: ", lines[x])
        table = lines[x].split()[5].split('.')[1]
        keyseq = lines[x].split()[5].split('.')[2][:-1]
        #print("table name: " + table) # table name
        #print("seq key: " + keyseq) # table name
        if table not in sequences:
            #print("adding " + table + " to sequences")
            sequences[table] = []
        sequences[table].append(keyseq)
            
#print("Primary Keys")
#print(keys)
#print("\nSequences")
#print(sequences)


s = ""
schema = ""
tablename = ""
cols = []

# clear query.sql before writing
file = open(path + outputfile, 'w')
file.write('')
file.close()

for x in range(len(lines)):
    #print("line", lines[x]);
    lines[x] = re.sub("[^a-zA-Z0-9' ();,._]+", '', lines[x])
    

    if len(lines[x].split()) > 0 and lines[x].split()[0] == "CREATE" and lines[x].split()[1] == "TABLE": # first line of table
        schema += lines[x] + "\n"
        tablename = lines[x].split()[2]
        t2 = tablename.replace('public.', '')
        if t2 not in sequences:
            sequences[t2] = []
        if t2 not in keys:
            keys[t2] = []
            #print(keys, sequences)
        
    elif lines[x] == ");" and tablename != "": # last lines line of table
        tablename = tablename.replace('public.', '')
        s += makequeries(tablename, cols, path + outputfile)
        cols = []
        tablename = ""
        schema += ");\n\n"
        continue
    elif len(lines[x].split()) > 1 and tablename != "":
        cols.append(lines[x].split()[0]); # add field name to cols
        #if cls len == append primary key to line
        schemaLine = lines[x]

        # else:
        #     schemaLine += ","
        schema += schemaLine + "\n"
        
    if (len(lines[x].split()) == 1 and lines[x] == ");"): # end of table test?
        schema += lines[x] + "\n"
        
        continue
    
schema = schema.replace("public.", '')

file = open(path + r"\schema.sql", 'w')
file.write(schema)
file.close()

additionalQueries(path + outputfile)
