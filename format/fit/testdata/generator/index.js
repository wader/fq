// import xlsx from 'node-xlsx';
var xlsx = require('node-xlsx').default;

const command = process.argv[2];
const profilePath = process.argv[3];

if (!command) {
   console.log("Usage: node index.js [command] [profilePath]");
   console.log("");
   console.log("t - generate TypeDefMap");
   console.log("m - generate MsgDefMap");
   return 0;
}

const workSheetsFromFile = xlsx.parse(profilePath);

const inTypes = workSheetsFromFile[0].data;
const inMessages = workSheetsFromFile[1].data;

let currentType = '';
let dynamicFieldName = '';
let currentFDefNo = null;
const outTypes = {};
const outMessages = {};
const outSubFields = {};

for (let li = 1; li < inTypes.length; li++) {
   const row = inTypes[li];

   if (row[0]) {
      currentType = row[0];
      outTypes[currentType] = { type: row[1], fields: [] };
   } else {
      if (row[4] && row[4].indexOf("Deprecated") > -1) {
         continue;
      }
      const val = row[3];
      outTypes[currentType].fields[val] = row[2];
   }

}
for (let li = 1; li < inMessages.length; li++) {
   const row = inMessages[li];
   const refFields = {}

   if (row[0]) {
      currentType = row[0];
      currentMsgNum = outTypes.mesg_num.fields.indexOf(currentType);
      outMessages[currentMsgNum] = { msgNum: currentMsgNum, type: currentType, fields: {} };
      outSubFields[currentMsgNum] = { msgNum: currentMsgNum, fields: {} }
   } else {
      if (row[2] == undefined) {
         continue;
      }
      const fDefNo = row[1];
      const name = row[2];
      const type = row[3];
      const scale = row[6];
      const offset = row[7];
      const unit = row[8];

      if (fDefNo != null) {
         outMessages[currentMsgNum].fields[name] = { fDefNo, name, type, unit, scale, offset };
         currentFDefNo = fDefNo
         dynamicFieldName = name
      } else {
         const refField = row[11].split(",")[0]
         const refVals = row[12].split(",")

         if (!Object.hasOwnProperty.call(refFields, refField)) {
            refFields[refField] = {}
         }
         refVals.forEach(element => {
            refFields[refField][element] = { name, type, unit, scale, offset }
         });

         outMessages[currentMsgNum].fields[dynamicFieldName]["hasSub"] = true;
         outSubFields[currentMsgNum].fields[currentFDefNo] = refFields
      }
   }

}


if (command == "t") {

   console.log("package mappers");
   console.log("");

   console.log("var TypeDefMap = map[string]typeDefMap{");
   for (const key in outTypes) {
      if (Object.hasOwnProperty.call(outTypes, key)) {
         const element = outTypes[key];
         console.log(`\t\"${key}\": {`);

         for (let index = 0; index < element.fields.length; index++) {
            const field = element.fields[index];
            if (field) {
               console.log(`\t\t${index}: {Name: \"${field}\"},`);
            }
         }

         console.log(`\t},`);
      }
   }
   console.log(`}`);
}

if (command == "m") {
   console.log("package mappers");
   console.log("");

   const baseTypes = ["bool", "byte", "enum", "uint8", "uint8z", "sint8", "sint16", "uint16", "uint16z", "sint32",
      "uint32", "uint32z", "float32", "float64", "sint64", "uint64", "uint64z", "string"];

   console.log("var SubFieldDefMap = map[uint64]map[uint64]map[string]map[string]FieldDef{");
   for (const key in outSubFields) {
      if (Object.hasOwnProperty.call(outSubFields, key)) {
         const element = outSubFields[key];
         if (Object.keys(element.fields).length == 0) {
            continue
         }
         console.log(`\t${key}: {`);

         for (const fieldKey in element.fields) {
            const field = element.fields[fieldKey];
            console.log(`\t\t${fieldKey}: {`);

            for (const refFieldKey in field) {
               const refField = field[refFieldKey];
               console.log(`\t\t\t"${refFieldKey}": {`);

               for (const refValKey in refField) {
                  const subField = refField[refValKey];

                  if (subField) {
                     let type = "";
                     let unit = "";
                     let scale = "";
                     let offset = "";

                     if (baseTypes.indexOf(subField.type) == -1) {
                        type = `, Type: \"${subField.type}\"`;
                     }

                     if (subField.unit) {
                        unit = `, Unit: \"${subField.unit}\"`;
                     }

                     if (subField.scale) {
                        if (typeof (subField.scale) == "string") {
                           // ignore multi scale (for component fields) for now
                           const testScale = subField.scale.split(",");
                           if (testScale.length == 1) {
                              scale = `, Scale: ${testScale[0]}`;
                           }
                        } else {
                           scale = `, Scale: ${subField.scale}`;
                        }
                     }

                     if (subField.offset) {
                        offset = `, Offset: ${subField.offset}`
                     }

                     console.log(`\t\t\t\t"${refValKey}": {Name: \"${subField.name}\"${type}${unit}${scale}${offset}},`);
                  }
               }
               console.log(`\t\t\t},`);
            }
            console.log(`\t\t},`);
         }
         console.log(`\t},`);
      }
   }
   console.log("}");
   console.log("");

   console.log("var FieldDefMap = map[uint64]fieldDefMap{");

   for (const key in outMessages) {
      if (Object.hasOwnProperty.call(outMessages, key)) {
         const element = outMessages[key];
         console.log(`\t${key}: {`);

         for (const msgKey in element.fields) {
            if (Object.hasOwnProperty.call(element.fields, msgKey)) {
               const field = element.fields[msgKey];
               if (field) {
                  let type = "";
                  let unit = "";
                  let scale = "";
                  let offset = "";
                  let hasSub = "";

                  if (baseTypes.indexOf(field.type) == -1) {
                     type = `, Type: \"${field.type}\"`;
                  }

                  if (field.unit) {
                     unit = `, Unit: \"${field.unit}\"`;
                  }

                  if (field.scale) {
                     if (typeof (field.scale) == "string") {
                        // ignore multi scale (for component fields) for now
                        const testScale = field.scale.split(",");
                        if (testScale.length == 1) {
                           scale = `, Scale: ${testScale[0]}`;
                        }
                     } else {
                        scale = `, Scale: ${field.scale}`;
                     }
                  }

                  if (field.offset) {
                     offset = `, Offset: ${field.offset}`
                  }

                  if (field.hasSub) {
                     hasSub = `, HasSubField: ${field.hasSub}`
                  }

                  console.log(`\t\t${field.fDefNo}: {Name: \"${field.name}\"${type}${unit}${scale}${offset}${hasSub}},`);

               }
            }
         }

         console.log(`\t},`);
      }
   }
   console.log(`}`);
}
