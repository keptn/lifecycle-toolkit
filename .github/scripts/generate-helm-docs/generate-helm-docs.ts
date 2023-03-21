import {ManifestPlugin} from "release-please/build/src/plugin";
import {CandidateReleasePullRequest} from "release-please/build/src/manifest";
import runReadmeGenerator from "@bitnami/readme-generator-for-helm";

export default class HelmDocsGenerator extends ManifestPlugin {
    async run(pullRequests: CandidateReleasePullRequest[]): Promise<CandidateReleasePullRequest[]> {
        runReadmeGenerator();
        return pullRequests;
    }
}
