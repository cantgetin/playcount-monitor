interface PlaysSummaryProps {
    data: UserStatsDataset[];
}

const PlaysSummary = (props: PlaysSummaryProps) => {
    return (
        <>
            <h1 className="text-2xl">{props.data[props.data.length - 1].map_count} maps fetched</h1>
            <div>
                <div className="text-xl text-yellow-200">{props.data[props.data.length - 1].play_count} Plays now</div>
                <div className="text-sm text-orange-200">{props.data[props.data.length - 2].play_count} plays last time</div>
            </div>
            <div className="flex flex-col mt-auto ml-auto">
                <div className="flex gap-2 justify-center items-center ml-auto px-2">
                    <h1 className="text-xs text-green-300">▲</h1>
                    <h1 className="text-2xl text-green-300">
                        {props.data[props.data.length - 1].play_count - props.data[props.data.length - 2].play_count}
                    </h1>
                </div>
                <div className="text-xs text-zinc-400">total plays for last 24 hours</div>
            </div>

        </>
    );
};

export default PlaysSummary;