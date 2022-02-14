import * as grpc_node from "@grpc/grpc-js"
import { ok, err, Result } from 'neverthrow';
import { 
    ApiContainerServiceClientNode, 
    ExecCommandArgs, 
    ExecCommandResponse,
    newExecCommandArgs,
    ServiceID
 } from "../..";
import { ServiceContextBackend } from "./service_context";

export class GrpcNodeServiceContextBackend implements ServiceContextBackend {
    private readonly client: ApiContainerServiceClientNode
    constructor(client: ApiContainerServiceClientNode) {
        this.client = client
    }

    public async execCommand(command: string[], serviceId: ServiceID): Promise<Result<[number, string], Error>> {
        const args: ExecCommandArgs = newExecCommandArgs(serviceId, command);

        const promiseExecCommand: Promise<Result<ExecCommandResponse, Error>> = new Promise((resolve, _unusedReject) => {
            this.client.execCommand(args, (error: grpc_node.ServiceError | null, response?: ExecCommandResponse) => {
                if (error === null) {
                    if (!response) {
                        resolve(err(new Error("No error was encountered but the response was still falsy; this should never happen")));
                    } else {
                        resolve(ok(response!));
                    }
                } else {
                    resolve(err(error));
                }
            })
        });
        const resultExecCommand: Result<ExecCommandResponse, Error> = await promiseExecCommand;
        if (resultExecCommand.isErr()) {
            return err(resultExecCommand.error);
        }
        const resp: ExecCommandResponse = resultExecCommand.value;

        return ok([resp.getExitCode(), resp.getLogOutput()]);
    }
}