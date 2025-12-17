import { PDFParse } from 'pdf-parse';
import { readFile } from 'node:fs/promises';

import readline from "readline";

const readPdf = async (path: string) => {
    const buffer = await readFile(path);

    const parser = new PDFParse({ data: buffer });
    const textResult = await parser.getText({disableNormalization: true});
    await parser.destroy();
    const questionsRaw = textResult.text.match(/\d\n[A-Z"“”](?:(?!\d\n[A-Z"“”])[\s\S])*?A\)[\s\S]*?(?=\d\n[A-Z"“”]|$)/g)
    questionsRaw?.splice(-10)

    const questionsObj = questionsRaw?.map((value, i) => {
        let question = value.split("\nA)")

        if (!question[1]) {
            question = value.split("\n(A)")
            question[1] = "\n(A)" + question[1]
        } else question[1] = "\nA)" + question[1]

        let answers = question[1].match(/\n\(?[“”"A-Z]\)[\s\S]*?(?=\n\(?[“”"A-Z]\)|$)/g) ?? []
        answers[3] = answers[3]?.match(/^[\s\S]*\./g)?.[0]!

        if (i === 79 && answers[3].includes("-- ")) answers[3] = answers[3].split("-- ")[0]!

        question[0] = question[0]?.match(/(?<=\d\n)[\s\S]*/g)?.[0]!
        question[0] = question[0].replaceAll("\n", " ")


        const formatedAnswers = answers.map((answer, i) => {
            return answer.replaceAll("\n", " ").replace(/\s*\(?[“”"A-Z]\)\s*/, "")
        })

        return {
            statement: question[0],
            answers: formatedAnswers
        }
    })

    questionsObj?.forEach((v,i) => {
        console.log(`Questão ${i + 1}: ${v.statement}\n`)
    })
}

const readAnswerKey = async (path: string) => {
    const buffer = await readFile(path);

    const parser = new PDFParse({ data: buffer });
    const result = await parser.getText();
    await parser.destroy();

    const rawAnswerKeys = result.text.match(/1 2 3 4 5 6 7(?:.*\r?\n){8}/)?.[0].split("\n").filter((line, i) => {
        return i % 2 !== 0
    })

    let answerKeys: string[] = []

    rawAnswerKeys?.map(value => {
        return value.split(" ")
    }).forEach(value => {
        answerKeys.push(...value)
    })

    console.log(answerKeys)
}

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

/*
rl.question("Digite o caminho do arquivo da prova: ", (resposta) => {
    readPdf(`/home/m3raak1/Downloads/oab3.pdf`)
    rl.close();
});
*/
rl.question("Digito o nome do arquivo do gabarito: ", (resposta) => {
    readAnswerKey(`/home/m3raak1/Downloads/goab.pdf`)
    rl.close();
});


