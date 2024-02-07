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
const outTypes = {};
const outMessages = {};
const outFormatters = {};

for (let li = 1; li < inTypes.length; li++) {
   const row = inTypes[li];

   if (row[0]) {
      currentType = row[0];
      outTypes[currentType] = { type: row[1], fields: [] };
      outFormatters[currentType] = { type: row[1] };
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

   if (row[0]) {
      currentType = row[0];
      currentMsgNum = outTypes.mesg_num.fields.indexOf(currentType);
      outMessages[currentMsgNum] = { msgNum: currentMsgNum, type: currentType, fields: {} };
   } else {
      if (row[1] == undefined) {
         continue;
      }
      const fDefNo = row[1];
      const name = row[2];
      const type = row[3];
      const scale = row[6];
      const offset = row[7];
      const unit = row[8];

      outMessages[currentMsgNum].fields[name] = { fDefNo, name, type, unit, scale, offset };
      outFormatters[type] = { ...outFormatters[type], unit, scale, offset };
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
   console.log("var FieldDefMap = map[uint64]fieldDefMap{");

   const baseTypes = ["bool", "byte", "enum", "uint8", "uint8z", "sint8", "sint16", "uint16", "uint16z", "sint32",
      "uint32", "uint32z", "float32", "float64", "sint64", "uint64", "uint64z", "string"];

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

                  if (baseTypes.indexOf(field.type) == -1) {
                     type = `, Type: \"${field.type}\"`;
                  }

                  if (field.unit) {
                     unit = `, Unit: \"${field.unit}\"`;
                  }

                  if (field.scale) {
                     scale = `, Scale: ${field.scale}`;
                  }

                  if (field.offset) {
                     offset = `, Offset: ${field.offset}`
                  }

                  console.log(`\t\t${field.fDefNo}: {Name: \"${field.name}\"${type}${unit}${scale}${offset}},`);

               }
            }
         }

         console.log(`\t},`);
      }
   }
   console.log(`}`);
}
