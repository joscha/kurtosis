import {createEnclave} from "../../test_helpers/enclave_setup";
import {
    DEFAULT_DRY_RUN,
    EMPTY_EXECUTE_PARAMS,
    IS_PARTITIONING_ENABLED,
    JEST_TIMEOUT_MS,
    VALID_MODULE_NO_MODULE_INPUT_TYPE_REL_PATH,
    VALID_MODULE_NO_MODULE_INPUT_TYPE_TEST_NAME
} from "./shared_constants";
import * as path from "path";
import log from "loglevel";
import {generateScriptOutput, readStreamContentUntilClosed} from "../../test_helpers/startosis_helpers";
import {err} from "neverthrow";

jest.setTimeout(JEST_TIMEOUT_MS)

test("Test valid startosis module with no module input type in types file", async () => {
    // ------------------------------------- ENGINE SETUP ----------------------------------------------
    const createEnclaveResult = await createEnclave(VALID_MODULE_NO_MODULE_INPUT_TYPE_TEST_NAME, IS_PARTITIONING_ENABLED)

    if (createEnclaveResult.isErr()) {
        throw createEnclaveResult.error
    }

    const {enclaveContext, stopEnclaveFunction} = createEnclaveResult.value

    try {
        // ------------------------------------- TEST SETUP ----------------------------------------------
        const moduleRootPath = path.join(__dirname, VALID_MODULE_NO_MODULE_INPUT_TYPE_REL_PATH)

        log.info(`Loading module at path '${moduleRootPath}'`)

        const outputStream = await enclaveContext.executeKurtosisModule(moduleRootPath, EMPTY_EXECUTE_PARAMS, DEFAULT_DRY_RUN)
        if (outputStream.isErr()) {
            throw err(new Error(`An error occurred execute startosis module '${moduleRootPath}'`));
        }
        const [interpretationError, validationErrors, executionError, instructions] = await readStreamContentUntilClosed(outputStream.value);

        const expectedScriptOutput = "Hello world!\n"

        expect(generateScriptOutput(instructions)).toEqual(expectedScriptOutput)

        expect(interpretationError).toBeUndefined()
        expect(validationErrors).toEqual([])
        expect(executionError).toBeUndefined()
    } finally {
        stopEnclaveFunction()
    }
})
